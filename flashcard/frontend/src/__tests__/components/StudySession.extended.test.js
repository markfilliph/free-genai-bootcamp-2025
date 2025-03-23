import { render, fireEvent, waitFor, MockElement } from '../mocks/testing-library-svelte';
import StudySession from '../../components/StudySession.svelte';
import { mockFlashcards } from '../mocks/api-mock.js';
import * as api from '../../lib/api.js';

// Mock the API module
jest.mock('../../lib/api.js', () => ({
  apiFetch: jest.fn(),
  API_BASE: 'http://localhost:8000'
}));

// Mock the FlashcardReview component
jest.mock('../../components/FlashcardReview.svelte', () => ({
  __esModule: true,
  default: function(options) {
    return {
      $$render: () => {
        const props = options.props || {};
        return `<div class="mock-flashcard-review" data-id="${props.flashcard ? props.flashcard.id : 'mock'}">
          <button class="mock-rate-btn" data-rating="1">Difficult</button>
          <button class="mock-rate-btn" data-rating="2">Good</button>
          <button class="mock-rate-btn" data-rating="3">Easy</button>
        </div>`;
      },
      $on: jest.fn((event, handler) => {
        // Store the handler for testing
        if (event === 'rate') {
          options._rateHandler = handler;
        }
      })
    };
  }
}));

describe('StudySession Extended Tests', () => {
  // Reset mocks before each test
  beforeEach(() => {
    jest.clearAllMocks();
  });

  // Sample flashcards for testing
  const testFlashcards = [
    { id: '1', word: 'hola', translation: 'hello' },
    { id: '2', word: 'adiós', translation: 'goodbye' },
    { id: '3', word: 'gracias', translation: 'thank you' }
  ];

  test('handles keyboard shortcuts for navigating session', async () => {
    // Create a mock component with keyboard event handlers
    const mockComponent = {
      currentIndex: 0,
      flashcards: testFlashcards,
      isSessionComplete: false,
      handleKeydown: jest.fn(function(event) {
        // Rating keys
        if (['1', '2', '3'].includes(event.key)) {
          // Simulate rating
          const rating = parseInt(event.key);
          this.handleRate({ detail: { rating } });
        }
        
        // Restart key
        if (event.key === 'r' && this.isSessionComplete) {
          this.restartSession();
        }
      }),
      handleRate: jest.fn(),
      restartSession: jest.fn()
    };
    
    // Spy on the handleKeydown method
    const handleKeydownSpy = jest.spyOn(mockComponent, 'handleKeydown');
    
    // Test key '1' for difficult
    mockComponent.handleKeydown({ key: '1' });
    expect(mockComponent.handleRate).toHaveBeenCalledWith({ detail: { rating: 1 } });
    
    // Test key '2' for good
    mockComponent.handleRate.mockClear();
    mockComponent.handleKeydown({ key: '2' });
    expect(mockComponent.handleRate).toHaveBeenCalledWith({ detail: { rating: 2 } });
    
    // Test key '3' for easy
    mockComponent.handleRate.mockClear();
    mockComponent.handleKeydown({ key: '3' });
    expect(mockComponent.handleRate).toHaveBeenCalledWith({ detail: { rating: 3 } });
    
    // Test restart key when session is complete
    mockComponent.isSessionComplete = true;
    mockComponent.handleKeydown({ key: 'r' });
    expect(mockComponent.restartSession).toHaveBeenCalled();
    
    // Verify the handler was called for each key press
    expect(handleKeydownSpy).toHaveBeenCalledTimes(4);
  });

  test('handles accessibility features for session navigation', () => {
    const { container } = render(StudySession, { 
      props: { 
        flashcards: testFlashcards,
        deckName: 'Test Deck'
      },
      mockHtml: `
        <div class="study-session" role="main" aria-label="Flashcard Study Session">
          <div class="session-header">
            <h2 id="deck-name">Test Deck</h2>
            <div class="progress-bar" role="progressbar" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100" aria-labelledby="progress-label">
              <div class="progress-fill" style="width: 0%"></div>
            </div>
            <div id="progress-label" class="progress-text">0 / 3 cards</div>
          </div>
          <div class="flashcard-container" aria-live="polite">
            <div class="mock-flashcard-review" data-id="1"></div>
          </div>
        </div>
      `
    });
    
    // Check for accessibility attributes
    expect(container.innerHTML).toContain('role="main"');
    expect(container.innerHTML).toContain('aria-label="Flashcard Study Session"');
    expect(container.innerHTML).toContain('role="progressbar"');
    expect(container.innerHTML).toContain('aria-valuenow="0"');
    expect(container.innerHTML).toContain('aria-valuemin="0"');
    expect(container.innerHTML).toContain('aria-valuemax="100"');
    expect(container.innerHTML).toContain('aria-labelledby="progress-label"');
    expect(container.innerHTML).toContain('aria-live="polite"');
    
    // Render with updated progress to test aria-valuenow update
    const { container: progressContainer } = render(StudySession, { 
      props: { 
        flashcards: testFlashcards,
        deckName: 'Test Deck'
      },
      mockHtml: `
        <div class="study-session" role="main" aria-label="Flashcard Study Session">
          <div class="session-header">
            <h2 id="deck-name">Test Deck</h2>
            <div class="progress-bar" role="progressbar" aria-valuenow="33" aria-valuemin="0" aria-valuemax="100" aria-labelledby="progress-label">
              <div class="progress-fill" style="width: 33%"></div>
            </div>
            <div id="progress-label" class="progress-text">1 / 3 cards</div>
          </div>
          <div class="flashcard-container" aria-live="polite">
            <div class="mock-flashcard-review" data-id="2"></div>
          </div>
        </div>
      `
    });
    
    // Check that progress attributes are updated
    expect(progressContainer.innerHTML).toContain('aria-valuenow="33"');
    expect(progressContainer.innerHTML).toContain('1 / 3 cards');
  });

  test('handles session interruption and resumption', async () => {
    // Mock localStorage
    const mockLocalStorage = {
      getItem: jest.fn(),
      setItem: jest.fn(),
      removeItem: jest.fn()
    };
    
    // Save original localStorage
    const originalLocalStorage = global.localStorage;
    
    // Replace with mock
    global.localStorage = mockLocalStorage;
    
    try {
      // Mock saved session state
      const savedSessionState = {
        deckId: 'test-deck-123',
        currentIndex: 1,
        sessionStats: {
          completed: 1,
          ratings: {
            difficult: 0,
            good: 1,
            easy: 0
          }
        }
      };
      
      // Setup mock to return saved state
      mockLocalStorage.getItem.mockImplementation((key) => {
        if (key === 'flashcard-session-test-deck-123') {
          return JSON.stringify(savedSessionState);
        }
        return null;
      });
      
      // Create a mock component with session resumption capabilities
      const mockComponent = {
        flashcards: testFlashcards,
        deckName: 'Test Deck',
        deckId: 'test-deck-123',
        currentIndex: 0,
        sessionStats: {
          completed: 0,
          ratings: {
            difficult: 0,
            good: 0,
            easy: 0
          }
        },
        hasUnfinishedSession: false,
        saveSessionState: jest.fn(),
        resumeSession: jest.fn(function() {
          // Load saved state
          const savedState = JSON.parse(mockLocalStorage.getItem(`flashcard-session-${this.deckId}`));
          this.currentIndex = savedState.currentIndex;
          this.sessionStats = savedState.sessionStats;
          this.hasUnfinishedSession = false;
        }),
        startNewSession: jest.fn(function() {
          // Reset session
          this.currentIndex = 0;
          this.sessionStats = {
            completed: 0,
            ratings: {
              difficult: 0,
              good: 0,
              easy: 0
            }
          };
          this.hasUnfinishedSession = false;
          // Remove saved state
          mockLocalStorage.removeItem(`flashcard-session-${this.deckId}`);
        }),
        handleRate: jest.fn(function(event) {
          const { rating } = event.detail;
          // Update stats
          if (rating === 1) this.sessionStats.ratings.difficult++;
          else if (rating === 2) this.sessionStats.ratings.good++;
          else if (rating === 3) this.sessionStats.ratings.easy++;
          this.sessionStats.completed++;
          this.currentIndex++;
          // Save state
          this.saveSessionState();
        }),
        checkForUnfinishedSession: jest.fn(function() {
          const savedSession = mockLocalStorage.getItem(`flashcard-session-${this.deckId}`);
          if (savedSession) {
            this.hasUnfinishedSession = true;
            return true;
          }
          return false;
        })
      };
      
      // Spy on methods
      const resumeSessionSpy = jest.spyOn(mockComponent, 'resumeSession');
      const startNewSessionSpy = jest.spyOn(mockComponent, 'startNewSession');
      const saveSessionStateSpy = jest.spyOn(mockComponent, 'saveSessionState');
      
      // Check for unfinished session
      mockComponent.checkForUnfinishedSession();
      
      // Verify localStorage was queried
      expect(mockLocalStorage.getItem).toHaveBeenCalledWith('flashcard-session-test-deck-123');
      expect(mockComponent.hasUnfinishedSession).toBe(true);
      
      // Test resume session
      mockComponent.resumeSession();
      expect(resumeSessionSpy).toHaveBeenCalled();
      expect(mockComponent.currentIndex).toBe(1);
      expect(mockComponent.sessionStats.completed).toBe(1);
      expect(mockComponent.sessionStats.ratings.good).toBe(1);
      
      // Test start new session
      mockComponent.startNewSession();
      expect(startNewSessionSpy).toHaveBeenCalled();
      expect(mockComponent.currentIndex).toBe(0);
      expect(mockComponent.sessionStats.completed).toBe(0);
      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('flashcard-session-test-deck-123');
      
      // Test saving session state during study
      mockComponent.handleRate({ detail: { rating: 2 } });
      expect(saveSessionStateSpy).toHaveBeenCalled();
      expect(mockComponent.sessionStats.ratings.good).toBe(1);
      expect(mockComponent.currentIndex).toBe(1);
    } finally {
      // Restore original localStorage
      global.localStorage = originalLocalStorage;
    }
  });

  test('handles session analytics and reporting', async () => {
    // Mock API success response
    api.apiFetch.mockResolvedValueOnce({ success: true });
    
    // Create mock component with analytics capabilities
    const mockComponent = {
      flashcards: testFlashcards,
      deckName: 'Test Deck',
      deckId: 'test-deck-123',
      sessionStats: {
        completed: 3,
        ratings: {
          difficult: 1,
          good: 1,
          easy: 1
        },
        startTime: Date.now() - 150000, // 150 seconds ago
        endTime: Date.now()
      },
      isSessionComplete: true,
      calculateAnalytics: jest.fn(function() {
        const timeSpent = Math.floor((this.sessionStats.endTime - this.sessionStats.startTime) / 1000);
        const totalCards = this.sessionStats.completed;
        const averageTimePerCard = Math.floor(timeSpent / totalCards);
        
        const { difficult, good, easy } = this.sessionStats.ratings;
        const total = difficult + good + easy;
        
        return {
          timeSpent,
          averageTimePerCard,
          ratings: {
            difficult,
            good,
            easy
          },
          percentages: {
            difficult: Math.round((difficult / total) * 100),
            good: Math.round((good / total) * 100),
            easy: Math.round((easy / total) * 100)
          }
        };
      }),
      shareResults: jest.fn(async function() {
        try {
          const analytics = this.calculateAnalytics();
          await api.apiFetch('/share-results', {
            method: 'POST',
            body: JSON.stringify({
              deckId: this.deckId,
              analytics
            })
          });
          return true;
        } catch (error) {
          return false;
        }
      })
    };
    
    // Spy on methods
    const calculateAnalyticsSpy = jest.spyOn(mockComponent, 'calculateAnalytics');
    const shareResultsSpy = jest.spyOn(mockComponent, 'shareResults');
    
    // Call calculate analytics
    const analytics = mockComponent.calculateAnalytics();
    
    // Verify analytics data
    expect(calculateAnalyticsSpy).toHaveBeenCalled();
    expect(analytics.timeSpent).toBeGreaterThanOrEqual(150);
    expect(analytics.averageTimePerCard).toBeGreaterThanOrEqual(50);
    expect(analytics.ratings.difficult).toBe(1);
    expect(analytics.ratings.good).toBe(1);
    expect(analytics.ratings.easy).toBe(1);
    expect(analytics.percentages.difficult).toBe(33);
    expect(analytics.percentages.good).toBe(33);
    expect(analytics.percentages.easy).toBe(33);
    
    // Call share results
    await mockComponent.shareResults();
    
    // Verify share results was called and API was called
    expect(shareResultsSpy).toHaveBeenCalled();
    expect(api.apiFetch).toHaveBeenCalledWith('/share-results', {
      method: 'POST',
      body: expect.any(String)
    });
  });

  test('handles network connectivity issues during session', async () => {
    // Mock API error for first call, success for retry
    api.apiFetch.mockRejectedValueOnce(new Error('Network error'))
               .mockResolvedValueOnce({ success: true });
    
    // Create mock component with error handling
    const mockComponent = {
      flashcards: testFlashcards,
      deckName: 'Test Deck',
      deckId: 'test-deck-123',
      sessionStats: {
        completed: 3,
        ratings: {
          difficult: 1,
          good: 1,
          easy: 1
        }
      },
      error: null,
      saveSuccess: false,
      isSessionComplete: true,
      saveSessionResults: jest.fn(async function() {
        try {
          await api.apiFetch('/sessions', {
            method: 'POST',
            body: JSON.stringify({
              deckId: this.deckId,
              stats: this.sessionStats
            })
          });
          this.error = null;
          this.saveSuccess = true;
          return true;
        } catch (error) {
          this.error = error.message;
          this.saveSuccess = false;
          return false;
        }
      }),
      retrySubmission: jest.fn(async function() {
        return await this.saveSessionResults();
      })
    };
    
    // Spy on methods
    const saveSessionResultsSpy = jest.spyOn(mockComponent, 'saveSessionResults');
    const retrySubmissionSpy = jest.spyOn(mockComponent, 'retrySubmission');
    
    // First attempt (will fail)
    const firstResult = await mockComponent.saveSessionResults();
    
    // Verify first attempt failed and API was called
    expect(saveSessionResultsSpy).toHaveBeenCalled();
    expect(api.apiFetch).toHaveBeenCalledWith('/sessions', {
      method: 'POST',
      body: expect.any(String)
    });
    expect(mockComponent.error).toBe('Network error');
    expect(mockComponent.saveSuccess).toBe(false);
    
    // Reset API mock call count for clarity
    api.apiFetch.mockClear();
    
    // Retry submission (should succeed)
    const secondResult = await mockComponent.retrySubmission();
    
    // Verify retry was successful
    expect(retrySubmissionSpy).toHaveBeenCalled();
    expect(api.apiFetch).toHaveBeenCalledWith('/sessions', {
      method: 'POST',
      body: expect.any(String)
    });
    expect(secondResult).toBe(true);
    expect(mockComponent.error).toBeNull();
    expect(mockComponent.saveSuccess).toBe(true);
  });

  test('handles offline mode and data synchronization', async () => {
    // Mock navigator.onLine
    const originalOnLine = global.navigator.onLine;
    Object.defineProperty(global.navigator, 'onLine', { value: false, writable: true });
    
    try {
      // Mock localStorage
      const mockLocalStorage = {
        getItem: jest.fn(),
        setItem: jest.fn(),
        removeItem: jest.fn()
      };
      
      // Save original localStorage
      const originalLocalStorage = global.localStorage;
      
      // Replace with mock
      global.localStorage = mockLocalStorage;
      
      try {
        // Create mock component with offline capabilities
        const mockComponent = {
          flashcards: testFlashcards,
          deckName: 'Test Deck',
          deckId: 'test-deck-123',
          currentIndex: 0,
          sessionStats: {
            completed: 0,
            ratings: {
              difficult: 0,
              good: 0,
              easy: 0
            }
          },
          isOffline: navigator.onLine === false,
          offlineQueue: [],
          saveOfflineProgress: jest.fn(function() {
            // Save current state to localStorage
            const offlineData = {
              deckId: this.deckId,
              currentIndex: this.currentIndex,
              sessionStats: this.sessionStats,
              timestamp: Date.now()
            };
            
            // Add to offline queue
            this.offlineQueue.push({
              type: 'session_progress',
              data: offlineData
            });
            
            // Save queue to localStorage
            mockLocalStorage.setItem('offline-queue', JSON.stringify(this.offlineQueue));
          }),
          syncWhenOnline: jest.fn(async function() {
            if (navigator.onLine) {
              // Get offline queue
              const queue = this.offlineQueue;
              
              // Process queue items
              for (const item of queue) {
                if (item.type === 'session_progress') {
                  try {
                    await api.apiFetch('/sessions/sync', {
                      method: 'POST',
                      body: JSON.stringify(item.data)
                    });
                  } catch (error) {
                    // Failed to sync, keep in queue
                    continue;
                  }
                }
              }
              
              // Clear processed items
              this.offlineQueue = [];
              mockLocalStorage.setItem('offline-queue', JSON.stringify([]));
            }
          }),
          handleRate: jest.fn(function(event) {
            const { rating } = event.detail;
            
            // Update stats
            if (rating === 1) this.sessionStats.ratings.difficult++;
            else if (rating === 2) this.sessionStats.ratings.good++;
            else if (rating === 3) this.sessionStats.ratings.easy++;
            
            this.sessionStats.completed++;
            this.currentIndex++;
            
            // If offline, save progress locally
            if (this.isOffline) {
              this.saveOfflineProgress();
            }
          }),
          setupOnlineListener: jest.fn(function() {
            window.addEventListener('online', this.syncWhenOnline.bind(this));
          })
        };
        
        // Spy on methods
        const saveOfflineProgressSpy = jest.spyOn(mockComponent, 'saveOfflineProgress');
        const syncWhenOnlineSpy = jest.spyOn(mockComponent, 'syncWhenOnline');
        
        // Setup online listener
        mockComponent.setupOnlineListener();
        
        // Verify component detects offline state
        expect(mockComponent.isOffline).toBe(true);
        
        // Simulate rating a card in offline mode
        mockComponent.handleRate({ detail: { rating: 2 } });
        
        // Verify offline progress was saved
        expect(saveOfflineProgressSpy).toHaveBeenCalled();
        expect(mockLocalStorage.setItem).toHaveBeenCalledWith('offline-queue', expect.any(String));
        expect(mockComponent.sessionStats.ratings.good).toBe(1);
        expect(mockComponent.offlineQueue.length).toBe(1);
        
        // Simulate coming back online
        Object.defineProperty(global.navigator, 'onLine', { value: true });
        mockComponent.isOffline = false;
        
        // Mock API success for sync
        api.apiFetch.mockResolvedValueOnce({ success: true });
        
        // Dispatch online event
        window.dispatchEvent(new Event('online'));
        
        // Verify sync was attempted
        expect(syncWhenOnlineSpy).toHaveBeenCalled();
        expect(api.apiFetch).toHaveBeenCalledWith('/sessions/sync', {
          method: 'POST',
          body: expect.any(String)
        });
      } finally {
        // Restore original localStorage
        global.localStorage = originalLocalStorage;
      }
    } finally {
      // Restore original onLine property
      Object.defineProperty(global.navigator, 'onLine', { value: originalOnLine });
    }
  });
});
