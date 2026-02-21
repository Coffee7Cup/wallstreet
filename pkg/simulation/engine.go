package simulation

import (
	"context"
	"sync"
	"time"

	"github.com/coffee7cup/wallstreet/pkg/db"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"go.uber.org/zap"
)

type Engine struct {
	store        *db.Store
	state        *models.SimulationState // authoritative state
	mu           sync.RWMutex
	tickDuration time.Duration
	stopChan     chan struct{}
	subscribers  map[chan models.SimulationBroadcast]struct{}
}

func NewEngine(store *db.Store) *Engine {
	// example: n tick per second (adjust as needed)
	tickInterval := time.Second * 3

	logs.Log.Info("Engine initialized", zap.Duration("tick_duration", tickInterval))

	engine := &Engine{
		store:        store,
		state:        &models.SimulationState{Tick: 0, IsActive: true, IsPaused: true},
		tickDuration: tickInterval,
		stopChan:     make(chan struct{}),
		subscribers:  make(map[chan models.SimulationBroadcast]struct{}),
	}

	// Load persisted state if exists
	persisted, err := store.LoadSimulationState(context.Background())
	if err != nil {
		logs.Log.Error("Failed to load persisted simulation state", zap.Error(err))
	} else if persisted != nil {
		engine.state = persisted
		logs.Log.Info("Restored simulation state from database",
			zap.Int("tick", engine.state.Tick),
			zap.Bool("is_active", engine.state.IsActive))
	}

	return engine
}

// TODO This fucntion i snot working properly
func (e *Engine) tick() {
	e.mu.RLock()
	prevTick := e.state.Tick
	e.mu.RUnlock()

	nextTick := prevTick + 1

	// Check if data exists for the next tick outside the lock
	stocks, err := e.store.GetStocksAtTick(context.Background(), nextTick)
	if err != nil || len(stocks) == 0 {
		logs.Log.Info("Simulation data ended or check failed, stopping engine", zap.Int("last_tick", prevTick), zap.Error(err))
		e.Stop()
		return
	}

	e.mu.Lock()
	// Double check if state hasn't been modified by another call (like Stop)
	if !e.state.IsActive || e.state.IsPaused {
		e.mu.Unlock()
		return
	}
	e.state.Tick = nextTick
	stateCopy := *e.state
	e.mu.Unlock()

	logs.Log.Debug("Engine ticking", zap.Int("new_tick", stateCopy.Tick), zap.Int("prev_tick", prevTick))

	// Persist state
	if err := e.store.SaveSimulationState(context.Background(), stateCopy); err != nil {
		logs.Log.Error("Failed to persist engine tick", zap.Error(err))
	}

	// Broadcast Tick state
	broadcastMsg := models.SimulationBroadcast{
		Tick:     stateCopy.Tick,
		IsActive: stateCopy.IsActive,
		IsPaused: stateCopy.IsPaused,
	}
	e.broadcast(broadcastMsg)
}

// Run is the main background loop of the engine.
func (e *Engine) Run(ctx context.Context) {
	ticker := time.NewTicker(e.tickDuration)
	defer ticker.Stop()

	e.state.IsActive = true
	e.state.IsPaused = false

	logs.Log.Info("Engine main loop started", zap.Duration("interval", e.tickDuration))
	for {
		select {
		case <-ticker.C:
			e.mu.Lock()
			if !e.state.IsActive || e.state.IsPaused {
				e.mu.Unlock()
				continue
			}
			e.mu.Unlock()
			e.tick()

		case <-e.stopChan:
			logs.Log.Info("Engine loop stopping via stopChan")
			return

		case <-ctx.Done():
			logs.Log.Info("Engine loop stopping via context", zap.Error(ctx.Err()))
			return
		}
	}
}

// Start initiates or restarts the simulation at the given tick.
func (e *Engine) Start(tick int) {
	e.mu.Lock()
	e.state.IsActive = true
	e.state.IsPaused = false
	e.state.Tick = tick
	stateCopy := *e.state
	e.mu.Unlock()

	logs.Log.Info("Engine started", zap.Int("start_tick", tick), zap.String("category", logs.CategoryEngine))
	_ = e.store.SaveSimulationState(context.Background(), stateCopy)
	e.broadcast(models.SimulationBroadcast{
		Tick:     stateCopy.Tick,
		IsActive: stateCopy.IsActive,
		IsPaused: stateCopy.IsPaused,
	})
}

// Stop halts the simulation permanently until Start is called again.
func (e *Engine) Stop() {
	e.mu.Lock()
	e.state.IsActive = false
	e.state.IsPaused = true // Ensure it stays stopped
	stateCopy := *e.state
	e.mu.Unlock()

	logs.Log.Info("Engine stopped", zap.String("category", logs.CategoryEngine))
	_ = e.store.SaveSimulationState(context.Background(), stateCopy)
	e.broadcast(models.SimulationBroadcast{
		Tick:     stateCopy.Tick,
		IsActive: false,
		IsPaused: true,
	})
}

// Pause suspends the simulation tick progression.
func (e *Engine) Pause() {
	e.mu.Lock()
	if !e.state.IsActive {
		e.mu.Unlock()
		logs.Log.Warn("Ignoring pause request: engine is stoped", zap.String("category", logs.CategoryEngine))
		return
	}
	e.state.IsPaused = true
	stateCopy := *e.state
	e.mu.Unlock()

	logs.Log.Info("Engine paused", zap.String("category", logs.CategoryEngine))
	_ = e.store.SaveSimulationState(context.Background(), stateCopy)
	e.broadcast(models.SimulationBroadcast{
		Tick:     stateCopy.Tick,
		IsActive: stateCopy.IsActive,
		IsPaused: true,
	})
}

// Resume continues the simulation tick progression.
func (e *Engine) Resume() {
	e.mu.Lock()
	if !e.state.IsActive {
		e.mu.Unlock()
		logs.Log.Warn("Ignoring resume request: engine is stopped", zap.String("category", logs.CategoryEngine))
		return
	}
	e.state.IsPaused = false
	stateCopy := *e.state
	e.mu.Unlock()

	logs.Log.Info("Engine resumed", zap.String("category", logs.CategoryEngine))
	_ = e.store.SaveSimulationState(context.Background(), stateCopy)
	e.broadcast(models.SimulationBroadcast{
		Tick:     stateCopy.Tick,
		IsActive: stateCopy.IsActive,
		IsPaused: false,
	})
}

func (e *Engine) GetState() models.SimulationState {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return *e.state
}

func (e *Engine) broadcast(msg models.SimulationBroadcast) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	logs.Log.Debug("Broadcasting state to internal subscribers", zap.Int("subscribers", len(e.subscribers)))
	for ch := range e.subscribers {
		select {
		case ch <- msg:
		default:
			logs.Log.Warn("Dropped broadcast for slow internal subscriber")
		}
	}
}

func (e *Engine) Subscribe() <-chan models.SimulationBroadcast {
	ch := make(chan models.SimulationBroadcast, 100)
	e.mu.Lock()
	e.subscribers[ch] = struct{}{}
	e.mu.Unlock()
	logs.Log.Debug("New internal subscriber added to engine")
	return ch
}
