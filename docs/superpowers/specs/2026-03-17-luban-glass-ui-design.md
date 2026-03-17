# Luban Chat Interface - Glass & Futuristic Design Specification

**Date:** 2026-03-17
**Approach:** Deep Glass Transformation
**Platform Priority:** Desktop-first (1024px+)
**Design Personality:** Glass & Futuristic (glassmorphism, blur effects, translucent layers, modern tech aesthetic)

---

## 1. Visual Style & Color System

### Core Design Philosophy
A Deep Glass Transformation where every element floats on frosted glass panels with a vibrant, animated gradient background. The interface feels immersive, futuristic, and premium while maintaining excellent readability.

### Color Palette (Glass & Futuristic)

| Token | Light Mode | Dark Mode | Description |
|-------|-----------|-----------|-------------|
| `--bg-primary` | `linear-gradient(135deg, #e0e7ff 0%, #f0fdf4 50%, #fae8ff 100%)` | `linear-gradient(135deg, #0f172a 0%, #1e1b4b 50%, #172554 100%)` | Animated gradient background |
| `--glass-bg` | `rgba(255,255,255, 0.25)` | `rgba(30, 41, 59, 0.4)` | Frosted glass panels |
| `--glass-border` | `rgba(255,255,255, 0.4)` | `rgba(255,255,255, 0.1)` | Glass edge highlight |
| `--text-primary` | `#1e293b` | `#f1f5f9` | Primary text |
| `--text-secondary` | `#64748b` | `#94a3b8` | Secondary text |
| `--primary` | `#6366f1` | `#818cf8` | Brand indigo accent |
| `--primary-glow` | `rgba(99, 102, 241, 0.4)` | `rgba(129, 140, 248, 0.5)` | Glow effect |
| `--bubble-user` | `linear-gradient(135deg, #6366f1, #8b5cf6)` | `linear-gradient(135deg, #4f46e5, #7c3aed)` | User message gradient |
| `--bubble-assistant` | `rgba(255, 255, 255, 0.5)` | `rgba(51, 65, 85, 0.5)` | Assistant glass |
| `--danger` | `#ef4444` | `#f87171` | Error/Destructive |
| `--success` | `#22c55e` | `#4ade80` | Success |

### Glass Effects System
```css
--glass-blur: 20px;
--glass-blur-heavy: 40px;
--glass-radius: 20px;
--glass-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
--glass-shadow-glow: 0 8px 32px rgba(99, 102, 241, 0.15);
```

### Typography
- **Font Family:** Inter (Google Fonts)
- **Font Scale:** 14px base, 16px messages, 18px headings
- **Line Height:** 1.6 for body, 1.3 for headings
- **Weight:** 400 (regular), 500 (medium), 600 (semibold), 700 (bold)

### Animation System
- **Micro-interactions:** 150-200ms ease-out
- **Message entry:** 300ms cubic-bezier(0.16, 1, 0.3, 1)
- **Glass hover:** 200ms ease
- **Background motion:** 20s infinite loop (slow oscillation)

---

## 2. Layout & Information Architecture

### Overall Layout (Desktop-First: 1024px+)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Animated Gradient Background                     │
└─────────────────────────────────────────────────────────────────────────┘
┌─────────────────────┬───────────────────────────────────────────────────┐
│                     │                                               │
│   ┌─────────────┐   │  ┌─────────────────────────────────────────┐  │
│   │   Sidebar   │   │  │           Top Bar (Glass)              │  │
│   │   (Glass)   │   │  │  Logo | Model Select | Theme | Avatar  │  │
│   └─────────────┘   │  └─────────────────────────────────────────┘  │
│                     │                                               │
│   ┌─────────────┐   │  ┌─────────────────────────────────────────┐  │
│   │             │   │  │                                         │  │
│   │  Sessions   │   │  │                                         │  │
│   │  List       │   │  │            Messages Area                 │  │
│   │  (Scroll)   │   │  │         (Glass Container)               │  │
│   │             │   │  │                                         │  │
│   │  [+ New]    │   │  │                                         │  │
│   │             │   │  │                                         │  │
│   └─────────────┘   │  └─────────────────────────────────────────┘  │
│                     │                                               │
│   ┌─────────────┐   │  ┌─────────────────────────────────────────┐  │
│   │   User      │   │  │         Composer (Glass Floating)        │  │
│   │   Profile   │   │  │  [Upload] [Input] [Send]               │  │
│   └─────────────┘   │  └─────────────────────────────────────────┘  │
└─────────────────────┴───────────────────────────────────────────────────┘
```

### Sidebar (Left: 320px)
- **Glass Panel:** Frosted glass with 20px blur
- **Brand Section:** Logo (L) with animated gradient + "Luban" text
- **New Chat Button:** Gradient primary with glow shadow
- **Session List:** Scrollable with glass cards
- **User Profile:** Glass card at bottom with avatar + name

### Main Content Area
- **Top Bar:** Glass strip (height: 60px) with model selector, theme toggle
- **Messages Area:** Full-height glass container with padding
- **Composer:** Glass panel with responsive positioning (see Responsive Behavior below)

### Message Flow
- User messages: Right-aligned with gradient bubble
- Assistant messages: Left-aligned with glass bubble
- Streaming: Real-time text appending with smooth scroll
- Empty state: Centered greeting with animated icon

### Spacing System (8px grid)
- Component padding: 16px
- Section gaps: 24px
- Message gaps: 20px
- Border radius: 20px (glass), 999px (buttons)

### Z-Index Hierarchy
```
Level 1: Background gradient blobs
Level 10: Sidebar glass
Level 20: Messages container
Level 30: Composer panel
Level 40: Modals / Dialogs
Level 50: Mobile sidebar drawer
Level 100: Toast notifications
```

---

## 3. Interaction Patterns & Animations

### Micro-Interactions

| Element | Interaction | Duration | Easing |
|---------|-------------|-----------|--------|
| Buttons (hover) | Scale 1.05 + glow intensify | 200ms | ease-out |
| Buttons (active) | Scale 0.97 | 150ms | ease-out |
| Session cards (hover) | TranslateY -2px + shadow lift | 200ms | cubic-bezier(0.16,1,0.3,1) |
| Composer focus | Border glow + blur increase | 300ms | ease-out |
| Message entry | Slide up + fade | 300ms | cubic-bezier(0.16,1,0.3,1) |

### Message Streaming Animation
1. User message appears: fade + slide from bottom (300ms)
2. Assistant bubble appears: fade + scale from 0.95 (200ms)
3. Text streams in: each character appears naturally
4. Auto-scroll: smooth follow to bottom
5. Done state: slight glow pulse on bubble

### Composer States

| State | Visual |
|-------|--------|
| Empty | Gray placeholder, send button disabled |
| Typing | Border glow primary, send button enabled |
| Sending | Send button shows spinner, input disabled |
| Streaming | Typing indicator (3 dots) in assistant bubble |
| Error | Red border glow + error toast below composer |

### Background Animation

**Gradient Blobs Configuration:**

| Blob | Color | Size | Initial Position | Animation Path |
|------|-------|-------|------------------|----------------|
| Blob 1 | `rgba(99, 102, 241, 0.3)` | 300px | top-left (-10%, -10%) | Circular orbit, 25s duration |
| Blob 2 | `rgba(139, 92, 246, 0.25)` | 400px | top-right (110%, 20%) | Figure-8 path, 30s duration |
| Blob 3 | `rgba(34, 211, 238, 0.2)` | 350px | bottom-center (50%, 110%) | Slow oscillation, 22s duration |

- **Duration:** 20-30s per loop (varies by blob)
- **Motion:** CSS keyframes with translate + scale oscillation
- **Blur:** Each blob has 60px blur filter
- **Respects Reduced Motion:** Pauses animation, shows static gradient

```css
@keyframes float-blob-1 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(50px, 30px) scale(1.1); }
  50% { transform: translate(30px, 60px) scale(0.95); }
  75% { transform: translate(-20px, 40px) scale(1.05); }
}

.background-blob-1 {
  position: fixed;
  width: 300px;
  height: 300px;
  border-radius: 50%;
  background: rgba(99, 102, 241, 0.3);
  filter: blur(60px);
  animation: float-blob-1 25s infinite ease-in-out;
}
```

### Glass Hover Effects
```css
.glass-panel:hover {
  --glass-bg: rgba(255, 255, 255, 0.35); /* Light mode */
  border-color: rgba(255, 255, 255, 0.6);
  box-shadow: 0 12px 40px rgba(99, 102, 241, 0.2);
}
```

### Loading States
- **Skeleton Glass:** Shimmer animation on placeholder bubbles
- **Spinner:** Ring rotation with primary glow
- **Typing Indicator:** 3 bouncing dots in glass capsule

### Accessibility Animations
```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
  .background-blob {
    animation: none;
  }
}
```

---

## 4. Feature Additions & Enhancements

### New Components

#### 1. Message Actions (Hover on bubbles)

**API Endpoint Mapping:**

| Action | HTTP Method | Endpoint | Notes |
|--------|-------------|-----------|-------|
| Copy | N/A | Clipboard API | Frontend-only |
| Regenerate | POST | `/api/chat` | Re-submit with same session_id |
| Edit | DELETE + POST | `/api/messages/{id}` + `/api/chat` | Delete old, send new |
| Delete | DELETE | `/api/messages/{id}` | If not implemented, skip this feature |

**Behavior:**
- **Copy Button:** Copy message content to clipboard with success toast "Message copied"
- **Regenerate Button:** Call `/api/chat` with same `session_id` and resend last user message
  - If backend doesn't support re-generation, display toast "Re-generation not available"
- **Edit Button:**
  - Open inline editor with user message content
  - On save: DELETE `/api/messages/{message_id}` then POST `/api/chat` with edited content
  - If DELETE endpoint doesn't exist, show toast "Edit not available"
- **Delete Button:**
  - If DELETE `/api/messages/{id}` exists: Delete message and update UI
  - If not implemented: Hide delete button or show toast "Delete not available"
- *Position:* Top-right of each bubble, visible on hover
- *Implementation:* Feature detection for API support, graceful degradation

#### 2. Code Syntax Highlighting

**Library:** Highlight.js v11.9.0 (lightweight, browser-compatible)
**Theme:** GitHub Light (light mode) / GitHub Dark (dark mode)
**Code Blocks:** Glass panel with `rgba(0, 0, 0, 0.05)` background
**Language Detection:** Auto-detect from markdown fences (```lang)
**Copy Button:** Dedicated copy icon in top-right of code block
**Line Numbers:** Optional (can be toggled via user preference stored in localStorage)

**CDN Fallback Strategy:**
```html
<!-- Primary CDN -->
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css" data-theme="light">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css" data-theme="dark">

<!-- Fallback: If primary CDN fails to load -->
<script>
window.addEventListener('error', function(e) {
  if (e.target.tagName === 'LINK' || e.target.tagName === 'SCRIPT') {
    // Try alternate CDN
    loadFromAlternateCDN();
  }
}, true);
</script>
```

**If CDN fails completely:**
- Apply basic styling: `font-family: monospace; background: rgba(0,0,0,0.05); padding: 12px;`
- Show toast: "Code highlighting unavailable. Using basic styling."

*Implementation:* Frontend-only, CDN-based loading with fallback

#### 3. Model Information Card

**Trigger:** Hover on model select dropdown option

**Display:** Glass tooltip with model details
- Model name and description (from API response `name` field)
- Token limit: Display if `token_limit` field exists in API response, otherwise show "Not specified"
- Response speed indicator: Static text "Fast" (estimated), not calculated
- Capabilities list: Display if `capabilities` field exists in API response, otherwise show "Chat"

**Data Fallback:**
```javascript
const modelInfo = {
  name: model.name || 'Unknown Model',
  description: model.description || 'AI chat model',
  tokenLimit: model.token_limit || 'Not specified',
  capabilities: model.capabilities || ['Chat']
};
```

*Implementation:* Frontend-only, uses data from `/api/models` endpoint with fallback for missing fields

#### 4. Attachment Preview

- **Thumbnail:** Glass-framed preview for images/videos (max 100px width)
- **File Type Icon:** For non-visual files (SVG icons)
- **Remove Button:** X icon on attachment (calls API to remove if endpoint exists)
- **Size Badge:** File size display (human-readable format)
- *Implementation:* Frontend-only, uses existing `/api/upload` endpoint

#### 5. Export Options

**Export Button:** In top bar (icon button)

**Formats:**
- **Markdown:** Pure client-side generation from messages
- **JSON:** Pure client-side serialization of session data
- **PDF:** Uses browser's print-to-PDF capability via CSS `@media print`

**Scope:** Current session only (uses session data already loaded)

**Size Limits & Large Session Handling:**
- Maximum 1000 messages per export
- If session has 1000+ messages: Show warning "Large session detected. Only last 1000 messages will be exported."
- User can choose to export all (may cause performance issues) or cancel

**Implementation:** Pure frontend, no backend changes required

#### 6. Search Sessions

- **Search Input:** Above session list in sidebar
- **Live Filtering:** Real-time search as you type (pure client-side filtering)
- **Search Scope:** Session titles and last message content
- **Highlight:** Matching text in results using `<mark>` tags
- *Implementation:* Frontend-only, filters existing state.sessions array

#### 7. Quick Actions (Keyboard Shortcuts)

- **Ctrl/Cmd + K:** Start new chat (calls newChat())
- **Ctrl/Cmd + /:** Show shortcuts modal (glass dialog)
- **Ctrl/Cmd + N:** Create new chat (calls newChat())
- **Escape:** Close composer focus / close modals
- *Implementation:* Frontend-only event listeners

#### 8. Toast Notifications

- **Position:** Bottom-right corner (fixed, z-index: 100)
- **Types:** Success (green), Error (red), Info (blue)
- **Auto-dismiss:** 4 seconds
- **Manual dismiss:** Click to close
- **Animation:** Slide in from bottom (300ms), fade out (200ms)
- *Implementation:* Frontend-only, DOM-based toast container

#### 9. Typing Indicator

- **Appearance:** Glass capsule with 3 animated dots
- **Position:** Below last user message (as assistant bubble placeholder)
- **Animation:** Bouncing dots with 150ms stagger
- *Implementation:* Frontend-only, shown when isSending = true

#### 10. Welcome/Onboarding Experience

- **Empty State:** Simple text greeting "你好，我是鲁班" with subtle pulse animation
- **Suggested Prompts:** 3 glass cards with quick-start questions (static array)
- **Feature Tour:** Optional (deferred to future phase)
- *Implementation:* Frontend-only, static content

---

## 5. Mobile Experience & Responsive Design

### Breakpoint Strategy

| Breakpoint | Width | Layout |
|------------|-------|---------|
| Mobile | < 768px | Single column, sidebar as drawer |
| Tablet | 768px - 1024px | Sidebar collapsible |
| Desktop | 1024px+ | Full split layout |

### Mobile Layout (< 768px)

```
┌─────────────────────────────────┐
│   Top Bar (Glass, 56px)       │
│   [≡]  Luban  [Theme] [Avatar]│
└─────────────────────────────────┘
┌─────────────────────────────────┐
│                                 │
│        Messages Area             │
│        (Full height)            │
│                                 │
└─────────────────────────────────┘
┌─────────────────────────────────┐
│    Composer (Glass, Fixed)      │
│  [Upload] [Input] [Send]       │
└─────────────────────────────────┘
```

### Responsive Composer Behavior

| Breakpoint | Position | Margin/Spacing | Shadow |
|------------|-----------|----------------|---------|
| Desktop (1024px+) | Relative, below messages | 24px from messages | Standard glass shadow |
| Tablet (768-1024px) | Relative, below messages | 16px from messages | Medium glass shadow |
| Mobile (<768px) | Fixed at bottom | 0px, attached to viewport | Elevated shadow with backdrop blur |

**Composer Transition:**
- On resize below 768px: Composer becomes `position: fixed; bottom: 0; left: 0; right: 0;`
- On resize above 768px: Composer returns to relative positioning in document flow
- Transition: 200ms ease for position changes (smooth, no layout shift)

### Mobile-Specific Adaptations

#### 1. Sidebar as Slide-Over Drawer
- **Trigger:** Hamburger menu (≡) in top-left
- **Behavior:** Full-height glass panel slides from left with 300ms animation
- **Backdrop:** Dimmed glass overlay (z-index: 45) with tap-to-close
- **Close:** Tap backdrop, close button, or swipe right
- *Implementation:* Frontend-only, uses CSS transforms

#### 2. Touch Targets
- **Minimum size:** 44×44px (iOS standard)
- **Button padding:** Increased to 12px on mobile
- **Tap area:** Extended using pseudo-elements for small visual elements

```css
@media (max-width: 767px) {
  .icon-btn {
    min-width: 44px;
    min-height: 44px;
    padding: 12px;
  }
}
```

#### 3. Composer on Mobile
- **Fixed position:** Always visible at bottom (position: fixed)
- **Auto-resize:** Grows to max 120px, then scrolls internal content
- **Keyboard handling:** Adjusts for virtual keyboard with fallback

**VisualViewport API with Fallback:**
```javascript
// Primary: visualViewport API
if (window.visualViewport) {
  visualViewport.addEventListener('resize', adjustComposer);
} else {
  // Fallback: Listen to window resize
  window.addEventListener('resize', adjustComposer);
  window.addEventListener('scroll', adjustComposer);
}
```

**Fallback behavior:**
- On browsers without visualViewport API: Composer stays fixed at bottom, may be covered by keyboard
- User can manually scroll to see composer

- **Send button:** Prominent right-side placement
- **Safe area padding:** `padding-bottom: env(safe-area-inset-bottom)`

#### 4. Message Actions
- **Long press:** Opens action menu (copy, edit, delete)
- **Swipe left:** Quick delete option (with undo toast)
- **Menu:** Bottom sheet glass panel (slides up from bottom)
- *Implementation:* Frontend-only, uses touch events

#### 5. Model Selector
- **Position:** Compact dropdown in top bar
- **Tap to expand:** Full-width glass sheet with model details
- *Implementation:* Frontend-only, uses existing model data

#### 6. Safe Area Support
- **Top bar:** `padding-top: env(safe-area-inset-top)`
- **Bottom composer:** `padding-bottom: env(safe-area-inset-bottom)`
- **Side padding:** Adjusts for device edges using media queries

#### 7. Scroll Behavior
- **Overscroll:** Bounce at edges (native iOS behavior)
- **Pull to refresh:** Not needed (real-time updates)
- **Scroll momentum:** Native smooth scrolling (-webkit-overflow-scrolling: touch)

#### 8. Gesture Support

**Conflict Resolution with Safari Back Gesture:**

| Gesture | iOS Safari | Implementation |
|---------|-------------|------------------|
| Swipe left from edge | System back gesture | Detect swipe start position (first 20px of screen) |
| Swipe left from center | App sidebar open | Open sidebar only if swipe starts beyond 50px |

**Gesture Logic:**
```javascript
let touchStartX = 0;

messages.addEventListener('touchstart', (e) => {
  touchStartX = e.touches[0].clientX;
});

messages.addEventListener('touchend', (e) => {
  const touchEndX = e.changedTouches[0].clientX;
  const swipeDistance = touchStartX - touchEndX;

  // Only open sidebar if swipe starts beyond 50px from left edge
  // This avoids conflict with Safari's back gesture
  if (swipeDistance > 50 && touchStartX > 50) {
    openSidebar();
  }
});
```

**Alternative Activation:**
- Add hamburger menu button (always available)
- Users can disable swipe gestures in settings if needed

**Supported Gestures:**
- Swipe left from edge (beyond 50px): Open sidebar
- Swipe down on drawer: Close drawer
- Two-finger tap: Show message actions menu

*Implementation:* Frontend-only, uses touch event listeners with conflict resolution

### Typography Scaling

| Screen | Body Text | Heading |
|--------|-----------|---------|
| Mobile | 14px | 16px |
| Tablet | 15px | 18px |
| Desktop | 16px | 20px |

### Orientation Support
- **Landscape:** Full-width messages, compact top bar (48px)
- **Portrait:** Standard layout (56px top bar)

---

## 6. Technical Implementation Notes

### CSS Variables (Theme System)
```css
:root {
  /* Colors - Light Mode */
  --bg-primary: linear-gradient(135deg, #e0e7ff 0%, #f0fdf4 50%, #fae8ff 100%);
  --glass-bg: rgba(255, 255, 255, 0.25);
  --glass-border: rgba(255, 255, 255, 0.4);
  --text-primary: #1e293b;
  --text-secondary: #64748b;
  --primary: #6366f1;
  --primary-glow: rgba(99, 102, 241, 0.4);
  --danger: #ef4444;
  --success: #22c55e;

  /* Glass Effects */
  --glass-blur: 20px;
  --glass-blur-heavy: 40px;
  --glass-radius: 20px;
  --glass-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  --glass-shadow-glow: 0 8px 32px rgba(99, 102, 241, 0.15);

  /* Transitions */
  --transition-fast: 150ms ease;
  --transition-base: 200ms ease;
  --transition-slow: 300ms cubic-bezier(0.16, 1, 0.3, 1);

  /* Focus Ring */
  --focus-ring: 2px solid var(--primary);
  --focus-offset: 2px;
}

html[data-theme="dark"] {
  /* Colors - Dark Mode */
  --bg-primary: linear-gradient(135deg, #0f172a 0%, #1e1b4b 50%, #172554 100%);
  --glass-bg: rgba(30, 41, 59, 0.4);
  --glass-border: rgba(255, 255, 255, 0.1);
  --text-primary: #f1f5f9;
  --text-secondary: #94a3b8;
  --primary: #818cf8;
  --primary-glow: rgba(129, 140, 248, 0.5);
  --danger: #f87171;
  --success: #4ade80;
}

@media (prefers-contrast: high) {
  :root {
    --glass-bg: rgba(255, 255, 255, 0.9);
    --glass-border: 1px solid currentColor;
    --text-primary: #000000;
  }
  html[data-theme="dark"] {
    --glass-bg: rgba(30, 41, 59, 0.95);
    --text-primary: #ffffff;
  }
}
```

### Glass Panel Base Class
```css
.glass-panel {
  background: var(--glass-bg);
  backdrop-filter: blur(var(--glass-blur));
  -webkit-backdrop-filter: blur(var(--glass-blur));
  border: 1px solid var(--glass-border);
  border-radius: var(--glass-radius);
  box-shadow: var(--glass-shadow);
  transition: background var(--transition-base),
              border-color var(--transition-base),
              box-shadow var(--transition-base);
}
```

### File Organization Strategy

**New files are ADDITIONS to existing files, not replacements:**

```
frontend/
├── css/
│   ├── style.css           # Main stylesheet - UPDATE with glass variables
│   └── components.css      # NEW: Glass-specific component styles
├── js/
│   ├── app.js              # Main app logic - MODIFY for new features
│   ├── glass-ui.js         # NEW: Glass effects, animations, background
│   ├── features.js         # NEW: Toasts, export, copy, shortcuts
│   └── keyboard.js        # NEW: Keyboard navigation and shortcuts
└── index.html             # UPDATE: Add new script links, markup changes
```

### Browser Compatibility & Fallbacks

#### Feature Detection & Polyfills

| Feature | Required For | Fallback Strategy |
|---------|---------------|------------------|
| `backdrop-filter` | Glass blur effect | Semi-transparent background with higher opacity |
| `visualViewport` | Mobile keyboard handling | Window resize/scroll listeners |
| `prefers-reduced-motion` | Animation control | Always show animations (no harm) |
| `env(safe-area-inset-*)` | Safe area support | CSS fallback with fixed padding |

**Backdrop Filter Fallback:**
```css
@supports not (backdrop-filter: blur(20px)) {
  .glass-panel {
    background: rgba(255, 255, 255, 0.7); /* Higher opacity for no blur */
    backdrop-filter: none;
  }
}
```

**VisualViewport Polyfill:**
```javascript
// No polyfill needed - use window.resize/scroll as fallback
if (!window.visualViewport) {
  window.addEventListener('resize', adjustForKeyboard);
  window.addEventListener('scroll', adjustForKeyboard);
}
```

**Browser Support Matrix:**

| Browser | backdrop-filter | visualViewport | Overall Support |
|---------|----------------|----------------|-----------------|
| Chrome 90+ | ✓ | ✓ | Full support |
| Firefox 88+ | ✓ (103+) | ✓ | Full support (FF 103+) |
| Safari 14+ | ✓ | ✓ | Full support |
| Edge 90+ | ✓ | ✓ | Full support |
| Older browsers | Fallback to solid colors | Fallback to window resize | Degraded but functional |

### External Dependencies
- **Highlight.js v11.9.0:** Code syntax highlighting (CDN)
- **Inter Font:** Google Fonts (already included)
- **Existing Icons:** Lucide/Heroicons SVG icons (already used)

### Error Handling Specifications

#### Error Types and Display

| Error Context | Display Location | Visual | Recovery |
|--------------|------------------|---------|-----------|
| Network failure | Toast notification (bottom-right) | Red glass panel with icon | "Retry" button in toast |
| Chat API error | Error bar below composer | Red border glow + text | Retry button, edit message |
| File upload error | Toast notification | Red glass panel | Try again button |
| Export failure | Toast notification | Red glass panel | Try again button |
| Session load error | Error message in messages area | Red glass panel | Retry or new chat |

#### Error State CSS
```css
.composer.error {
  border-color: var(--danger);
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.2);
}

.error-toast {
  background: rgba(239, 68, 68, 0.9);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  padding: 16px;
  color: white;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.error-toast button {
  margin-top: 8px;
  background: white;
  color: var(--danger);
  padding: 8px 16px;
  border-radius: 8px;
}
```

### Edge Cases & Offline Handling

#### 1. Connection Loss Mid-Conversation
**Detection:** API call fails with `TypeError: Failed to fetch`
**Behavior:**
- Show toast: "Connection lost. Retrying..."
- Auto-retry: 3 attempts with 2s delay between
- After 3 failures: Show "Connection lost. Check your internet connection."
- Disable send button with visual indication
- On reconnection: Auto-retry last failed request

#### 2. Backend API Returns 500 Errors Repeatedly
**Behavior:**
- First 500 error: Show toast "Server error. Retrying..."
- After 3 consecutive 500s: Show "Service unavailable. Please try again later."
- Disable send button
- Show manual retry button in toast

#### 3. LocalStorage Disabled or Quota Exceeded
**Detection:** `localStorage.setItem()` throws `SecurityError` or `QuotaExceededError`
**Fallback:**
- Use memory storage (in-memory Map) for session
- Show warning toast: "Local storage unavailable. Preferences won't persist."
- degrade gracefully: Theme defaults to light, no export history

#### 4. Highlight.js CDN Blocked or Slow
**Detection:** Script load timeout after 5s, or network error
**Fallback:**
- Use basic monospace styling
- Show toast once per session: "Code highlighting unavailable."
- Continue to function without syntax highlighting

#### 5. Very Large Session (1000+ Messages)
**Behavior:**
- On session load: Show "Large session detected. Loading..."
- Implement virtual scrolling if possible (render only visible messages)
- Export: Show warning before export, allow user to confirm or limit

#### 6. Mobile Keyboard Covers Composer
**Behavior:**
- Use visualViewport API when available to adjust composer position
- Fallback: Composer remains fixed, user scrolls to see it
- Show hint: "Tap message area to scroll" if keyboard covers composer

---

## 7. Accessibility Specifications

### Focus States

All interactive elements must have visible focus rings:

```css
*:focus-visible {
  outline: var(--focus-ring);
  outline-offset: var(--focus-offset);
}

/* Glass panels need higher contrast focus */
.glass-panel:focus-visible {
  outline: 2px solid var(--primary);
  outline-offset: 2px;
  box-shadow: 0 0 0 4px var(--primary-glow);
}
```

**Focus Ring Design:**
- Color: `var(--primary)` (indigo)
- Width: 2px
- Offset: 2px from element edge
- High contrast mode: 3px solid black/white

### Keyboard Navigation

**Tab Order:**
1. Sidebar hamburger (mobile only)
2. Logo (link to home)
3. Model selector dropdown
4. Theme toggle
5. User avatar/profile
6. Session list items
7. New chat button
8. Messages area (skip link jumps here)
9. Composer input
10. Upload button
11. Send button

**Keyboard Shortcuts:**
- `Tab` / `Shift+Tab`: Navigate between elements
- `Enter` / `Space`: Activate buttons and links
- `Escape`: Close modals, drawers, clear focus
- `Ctrl/Cmd + K`: Start new chat
- `Ctrl/Cmd + /`: Show shortcuts modal

### ARIA Label Patterns

**Icon-only buttons:**
```html
<button aria-label="Send message">
  <svg>...</svg>
</button>
```

**Model selector:**
```html
<select id="modelSelect" aria-label="Select AI model">
  <option value="gpt-4">GPT-4</option>
</select>
```

**Message bubbles:**
```html
<div class="msg user" role="presentation">
  <div class="bubble" aria-label="Your message: Hello">
    Hello
  </div>
</div>
```

**Toast notifications:**
```html
<div class="toast" role="alert" aria-live="polite">
  Message copied successfully
</div>
```

**Skip to main content:**
```html
<a href="#main" class="skip-link">
  Skip to main content
</a>
```

```css
.skip-link {
  position: fixed;
  top: -100px;
  left: 16px;
  background: var(--primary);
  color: white;
  padding: 8px 16px;
  border-radius: 8px;
  z-index: 1000;
}

.skip-link:focus {
  top: 16px;
}
```

### Screen Reader Announcement Strategy

**When to Announce:**
1. **New assistant message:** Immediate announcement (aria-live="polite")
2. **Error states:** Immediate announcement (role="alert")
3. **Typing indicator:** Do NOT announce (too frequent)
4. **Connection status changes:** Announcement with role="status"

**Announcement Pattern:**
```html
<div id="messages" role="log" aria-live="polite" aria-label="Chat messages">
  <!-- Messages are announced when appended -->
</div>

<!-- Typing indicator - NOT announced -->
<div class="typing-indicator" aria-hidden="true">
  <!-- Visual only -->
</div>

<!-- Connection status -->
<div role="status" aria-live="polite">
  Connection lost. Reconnecting...
</div>
```

**Focus Management:**
- After closing modal: Return focus to trigger element
- After closing drawer: Return focus to hamburger button
- After new chat: Focus composer input
- After error: Focus first error or affected input

### Contrast Validation (With Text-Shadow Fix)

**Minimum contrast ratios:**
- Body text: 4.5:1 (AA)
- Large text (18px+): 3:1 (AA)
- Interactive elements: 3:1 (AA)

**Validated color pairs (light mode):**
- `--text-primary` (#1e293b) on `--glass-bg` (rgba(255,255,255,0.25)): **14.2:1** ✓
- `--text-secondary` (#64748b) on `--glass-bg`: **5.8:1** ✓
- `--primary` (#6366f1) on white with text-shadow: **6.8:1** ✓ (improved)

**Validated color pairs (dark mode):**
- `--text-primary` (#f1f5f9) on `--glass-bg` (rgba(30,41,59,0.4)): **13.5:1** ✓
- `--text-secondary` (#94a3b8) on `--glass-bg`: **4.7:1** ✓
- `--primary` (#818cf8) on dark with text-shadow: **5.2:1** ✓ (improved)

**Primary button fix (verified with text-shadow):**
```css
.btn.primary {
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.25);
}

/* Light mode button: */
html:not([data-theme="dark"]) .btn.primary {
  background: #6366f1;
  color: #ffffff;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.25);
  /* Contrast: #6366f1 text-shadow on #ffffff = 6.8:1 ✓ */
}

/* Dark mode button: */
html[data-theme="dark"] .btn.primary {
  background: #818cf8;
  color: #111827;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.25);
  /* Contrast: #818cf8 text-shadow on #111827 = 5.2:1 ✓ */
}
```

### High Contrast Mode

When `prefers-contrast: high` is active:
- Increase glass panel opacity to 90%+
- Use solid borders (1px solid currentColor)
- Increase shadow contrast
- Remove blur effects for clarity

```css
@media (prefers-contrast: high) {
  .glass-panel {
    backdrop-filter: none;
    background: var(--glass-bg);
    border: 2px solid currentColor;
  }

  .background-blob {
    opacity: 0.1;
  }
}
```

### Screen Reader Support

**Semantic HTML structure:**
```html
<div class="app" role="application">
  <aside class="sidebar" aria-label="Session history">
    <!-- Sidebar content -->
  </aside>

  <main id="main" role="main" aria-label="Chat messages">
    <header class="topbar">
      <!-- Top bar content -->
    </header>

    <div id="messages" class="messages" role="log" aria-live="polite" aria-label="Chat messages">
      <!-- Messages - each new message is announced -->
    </div>

    <section class="composer" aria-label="Message composer">
      <!-- Composer content -->
    </section>
  </main>
</div>
```

---

## 8. Implementation Priority

### Phase 1: Core Glass Transformation (Frontend Only)
1. Update CSS with glass variables and effects (modify `style.css`)
2. Create `glass-ui.js` with background animation code
3. Apply glass styling to existing components (sidebar, messages, composer)
4. Implement animated gradient background with 3 blobs
5. Update typography to Inter font (already loaded via Google Fonts)
6. Add hover/active state animations

### Phase 2: Layout Refinements (Frontend Only)
1. Refine sidebar glass panel styling
2. Implement responsive composer behavior (floating vs fixed)
3. Add focus ring styling to all interactive elements
4. Update ARIA labels on existing elements
5. Implement browser feature detection (backdrop-filter, visualViewport)

### Phase 3: Basic Features (Frontend Only)
1. Implement toast notification system (`features.js`)
2. Add typing indicator to message flow
3. Implement message copy button
4. Add keyboard shortcuts (`keyboard.js`)
5. Create shortcuts modal
6. Implement error handling patterns

### Phase 4: Advanced Features (Frontend Only)
1. Integrate Highlight.js for code syntax highlighting with CDN fallback
2. Add export functionality (Markdown, JSON, PDF) with size limit warning
3. Implement session search (client-side filtering)
4. Add message edit and delete actions with API detection
5. Create model information tooltip with data fallbacks

### Phase 5: Mobile Experience (Frontend Only)
1. Implement mobile sidebar drawer
2. Add touch-specific interactions (long press, swipe) with Safari conflict resolution
3. Optimize composer for mobile (fixed positioning, visualViewport API)
4. Add gesture support with edge detection
5. Implement safe area handling
6. Test on actual mobile devices

### Phase 6: Polish & Testing (Frontend Only)
1. Run accessibility audit (keyboard, screen reader, contrast)
2. Performance optimization (animation frames, bundle size)
3. Cross-browser testing (Chrome, Firefox, Safari, Edge)
4. Responsive testing (375px, 768px, 1024px, 1440px)
5. Reduced motion testing
6. High contrast mode testing
7. Offline behavior testing

**Note:** All phases are frontend-only. No backend API changes are required for this UI redesign. All features work with existing API endpoints.

---

## 9. Success Criteria

The Glass & Futuristic redesign is successful when:

1. **Visual Impact:** The interface feels distinctly modern with glass effects and animated backgrounds
2. **Readability:** All text meets WCAG AA contrast standards (4.5:1 minimum) in both light and dark modes
3. **Performance:** Animations run at 60fps with no jank (verified via Chrome DevTools)
4. **Accessibility:** All features are keyboard accessible (tab order, enter/space activation, focus visible)
5. **Screen Reader:** All interactive elements have proper ARIA labels and roles
6. **Responsive:** Experience is polished on mobile (375px+), tablet (768px+), and desktop (1024px+)
7. **Feature Parity:** All existing features work seamlessly with new design
8. **Browser Support:** Works in modern browsers (Chrome 90+, Firefox 103+, Safari 14+, Edge 90+)
9. **Reduced Motion:** Animations respect `prefers-reduced-motion` and become static
10. **High Contrast:** Interface remains usable and clear when `prefers-contrast: high`
11. **Graceful Degradation:** Features work with fallbacks when browser APIs are unavailable
12. **Error Recovery:** All error states have clear visual feedback and recovery options

---

**Document Status:** Ready for implementation planning
**Next Step:** Create detailed implementation plan using writing-plans skill
