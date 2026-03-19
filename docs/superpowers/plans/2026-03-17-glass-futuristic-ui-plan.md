# Glass & Futuristic Luban UI Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Transform Luban chat interface into a Glass & Futuristic design with animated gradient backgrounds, frosted glass panels, and modern interactions while maintaining all existing functionality.

**Architecture:** Frontend-only transformation - no backend API changes required. All new features work with existing API endpoints with graceful degradation for missing capabilities. The design uses CSS custom properties (variables) for theming and glass effects that can be applied progressively.

**Tech Stack:** HTML5, CSS3 (backdrop-filter, CSS custom properties), JavaScript (ES6+), Highlight.js v11.9.0 (CDN-based code highlighting), Inter font (Google Fonts)

---

## File Structure

### Files to Modify

| File | Responsibility | Changes |
|------|----------------|----------|
| `frontend/css/style.css` | Main stylesheet | Replace CSS variables with glass system, add glass panel classes, update focus styles, add animations |
| `frontend/js/app.js` | Main app logic | Integrate new JS modules, update element references, add feature toggles |
| `frontend/index.html` | HTML structure | Add new script tags, update class names, add semantic markup |

### Files to Create

| File | Responsibility |
|------|----------------|
| `frontend/css/components.css` | Glass-specific component styles, animations, mobile adaptations |
| `frontend/js/glass-ui.js` | Background animation blob system, glass effect utilities |
| `frontend/js/features.js` | Toast notifications, export functionality, message actions, copy to clipboard |
| `frontend/js/keyboard.js` | Keyboard shortcuts, keyboard navigation, focus management |

---

## Chunk 1: CSS Foundation & Glass System

### Task 1: Update CSS Variables

**Files:**
- Modify: `frontend/css/style.css:1-31`

- [ ] **Step 1: Replace :root variables with glass system**

Replace lines 1-31 with:

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
  --bg-sidebar: rgba(255, 255, 255, 0.15);

  /* Message Bubble Colors */
  --bubble-user: linear-gradient(135deg, #6366f1, #8b5cf6);
  --bubble-assistant: rgba(255, 255, 255, 0.5);

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
  --bg-sidebar: rgba(30, 41, 59, 0.3);

  /* Message Bubble Colors - Dark Mode */
  --bubble-user: linear-gradient(135deg, #4f46e5, #7c3aed);
  --bubble-assistant: rgba(51, 65, 85, 0.5);
}
```

- [ ] **Step 2: Run browser test**

Open `frontend/index.html` in browser
Expected: Background shows gradient, UI renders with new colors

- [ ] **Step 3: Test dark mode toggle**

Click theme toggle button
Expected: Colors switch to dark mode palette

- [ ] **Step 4: Commit**

```bash
git add frontend/css/style.css
git commit -m "feat: replace CSS variables with glass system including bubble colors"
```

---

### Task 2: Add Glass Panel Base Class

**Files:**
- Modify: `frontend/css/style.css:900-958` (append to end)

- [ ] **Step 1: Add glass panel class at end of file**

```css
/* Glass Panel Base Class */
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

/* Glass Hover Effects */
.glass-panel:hover {
  --glass-bg: rgba(255, 255, 255, 0.35); /* Light mode */
  border-color: rgba(255, 255, 255, 0.6);
  box-shadow: 0 12px 40px rgba(99, 102, 241, 0.2);
}

html[data-theme="dark"] .glass-panel:hover {
  --glass-bg: rgba(30, 41, 59, 0.5);
  border-color: rgba(255, 255, 255, 0.15);
}
```

- [ ] **Step 2: Add browser compatibility fallback**

```css
@supports not (backdrop-filter: blur(20px)) {
  .glass-panel {
    background: rgba(255, 255, 255, 0.7);
    backdrop-filter: none;
  }
}
```

- [ ] **Step 3: Run browser test**

Open `frontend/index.html` in browser
Expected: No errors, styles load correctly

- [ ] **Step 4: Test in older browser**

Open in browser without backdrop-filter support
Expected: Falls back to solid background

- [ ] **Step 5: Commit**

```bash
git add frontend/css/style.css
git commit -m "feat: add glass panel base class with browser fallback"
```

---

### Task 3: Update Focus Styles

**Files:**
- Modify: `frontend/css/style.css:64-68`

- [ ] **Step 1: Replace focus-visible styles**

```css
/* Focus Styles - Accessibility */
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

- [ ] **Step 2: Test keyboard navigation**

Tab through interface, verify focus rings visible
Expected: All interactive elements show visible focus

- [ ] **Step 3: Test on glass panels**

Click sidebar, tab to next element
Expected: Glass panel focus has enhanced glow

- [ ] **Step 4: Commit**

```bash
git add frontend/css/style.css
git commit -m "feat: update focus styles with glass panel enhancement"
```

---

### Task 4: Add High Contrast Mode Support

**Files:**
- Modify: `frontend/css/style.css` (append after dark mode)

- [ ] **Step 1: Add high contrast media query**

```css
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

- [ ] **Step 2: Test high contrast mode**

Enable high contrast in OS settings, reload page
Expected: High opacity backgrounds, solid borders

- [ ] **Step 3: Test in both light and dark mode**

Toggle theme while in high contrast mode
Expected: Both themes maintain high contrast

- [ ] **Step 4: Commit**

```bash
git add frontend/css/style.css
git commit -m "feat: add high contrast mode support"
```

---

### Task 5: Add Reduced Motion Support

**Files:**
- Modify: `frontend/css/style.css` (append at end)

- [ ] **Step 1: Add reduced motion media query**

```css
/* Accessibility: Reduced Motion */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```

- [ ] **Step 2: Test reduced motion**

Enable reduced motion in OS settings, reload page
Expected: No animations, instant transitions

- [ ] **Step 3: Verify animations work normally**

Disable reduced motion, reload page
Expected: All animations play smoothly

- [ ] **Step 4: Commit**

```bash
git add frontend/css/style.css
git commit -m "feat: add reduced motion support"
```

---

## Chunk 2: Component Glass Styling

### Task 6: Apply Glass to Sidebar

**Files:**
- Modify: `frontend/css/style.css:86-93` (sidebar class)

- [ ] **Step 1: Update sidebar class**

```css
.sidebar {
  background: var(--glass-bg);
  backdrop-filter: blur(var(--glass-blur));
  -webkit-backdrop-filter: blur(var(--glass-blur));
  border-right: 1px solid var(--glass-border);
  display: flex;
  flex-direction: column;
  padding: 16px 12px;
  gap: 16px;
  transition: transform 300ms cubic-bezier(0.16, 1, 0.3, 1);
}

/* Sidebar hover effect - lift up */
.sidebar:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 24px var(--glass-shadow-glow);
}

@media (prefers-reduced-motion: reduce) {
  .sidebar:hover {
    transform: none;
  }
}
```

- [ ] **Step 2: Test sidebar rendering**

Open app, verify sidebar shows glass effect
Expected: Translucent background with blur

- [ ] **Step 3: Test in dark mode**

Toggle to dark theme
Expected: Sidebar glass effect persists

- [ ] **Step 4: Commit**

```bash
git add frontend/css/style.css
git commit -m "style: apply glass effect to sidebar with hover lift"
```

---

### Task 7: Apply Glass to Session Cards

**Files:**
- Modify: `frontend/css/style.css:153-190` (session-item class)

- [ ] **Step 1: Update session-item class**

```css
.session-item {
  background: var(--glass-bg);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid var(--glass-border);
  border-radius: 12px;
  padding: 10px 10px;
  display: grid;
  gap: 4px;
  cursor: pointer;
  transition: all var(--transition-base);
}

/* Session card hover effect - lift and glow */
.session-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 24px var(--glass-shadow-glow);
  border-color: rgba(255, 255, 255, 0.6);
}

.session-item.active {
  outline: 2px solid var(--primary);
  outline-offset: 1px;
  box-shadow: 0 0 0 4px var(--primary-glow);
  border-color: rgba(99, 102, 241, 0.3);
}

.session-item:active {
  transform: translateY(0);
  box-shadow: var(--glass-shadow);
}

@media (prefers-reduced-motion: reduce) {
  .session-item:hover {
    transform: none;
  }
}
```

- [ ] **Step 2: Test session hover**

Hover over session in sidebar
Expected: Smooth hover animation with lift

- [ ] **Step 3: Test active state**

Click a session to select it
Expected: Active session has visual indication

- [ ] **Step 4: Commit**

```bash
git add frontend/css/style.css
git commit -m "style: apply glass effect to session cards with hover effects"
```

---

### Task 8: Apply Glass to Message Bubbles

**Files:**
- Modify: `frontend/css/style.css:407-452` (bubble class)

- [ ] **Step 1: Update bubble classes**

```css
.bubble {
  max-width: min(760px, 92%);
  padding: 14px 18px;
  border-radius: 18px;
  border: 1px solid var(--glass-border);
  white-space: pre-wrap;
  line-height: 1.5;
  font-size: 15px;
  backdrop-filter: blur(var(--glass-blur));
  -webkit-backdrop-filter: blur(var(--glass-blur));
  transition: box-shadow var(--transition-base), transform var(--transition-base);
}

/* Message bubble hover - subtle shadow increase */
.bubble:hover {
  box-shadow: 0 12px 24px var(--glass-shadow-glow);
}

.msg.user .bubble {
  background: var(--bubble-user);
  color: #ffffff;
  margin-left: auto;
  border-radius: 18px 18px 4px 18px;
  border: none;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

.msg.assistant .bubble {
  background: var(--bubble-assistant);
  color: var(--text-primary);
  margin-right: auto;
  border-radius: 18px 18px 18px 4px;
}
```

- [ ] **Step 2: Test message rendering**

Open a chat with messages
Expected: Messages show glass/gradient effects using CSS variables

- [ ] **Step 3: Test bubble hover**

Hover over message bubble
Expected: Subtle shadow increase

- [ ] **Step 4: Commit**

```bash
git add frontend/css/style.css
git commit -m "style: apply glass effect to message bubbles using CSS variables"
```

---

### Task 9: Apply Glass to Top Bar

**Files:**
- Modify: `frontend/css/style.css:293-300` (topbar class)

- [ ] **Step 1: Update topbar class**

```css
.topbar {
  background: var(--glass-bg);
  backdrop-filter: blur(var(--glass-blur));
  -webkit-backdrop-filter: blur(var(--glass-blur));
  border-bottom: 1px solid var(--glass-border);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 16px 18px;
}

/* Top bar hover effect */
.topbar:hover {
  box-shadow: 0 4px 12px var(--glass-shadow-glow);
}

@media (prefers-reduced-motion: reduce) {
  .topbar:hover {
    box-shadow: none;
  }
}
```

- [ ] **Step 2: Test top bar rendering**

Open app, verify top bar glass effect
Expected: Translucent bar with blur

- [ ] **Step 3: Test scroll behavior**

Scroll messages past top bar
Expected: Content scrolls behind top bar

- [ ] **Step 4: Commit**

```bash
git add frontend/css/style.css
git commit -m "style: apply glass effect to top bar"
```

---

### Task 10: Apply Glass to Composer

**Files:**
- Modify: `frontend/css/style.css:462-476` (composer class)

- [ ] **Step 1: Update composer class**

```css
.composer {
  background: linear-gradient(
    to bottom,
    rgba(255, 255, 255, 0),
    rgba(255, 255, 255, 0.1)
  );
  border-top: 1px solid var(--glass-border);
  padding: 12px 18px 14px;
}

html[data-theme="dark"] .composer {
  background: linear-gradient(
    to bottom,
    rgba(0, 0, 0, 0),
    rgba(0, 0, 0, 0.1)
  );
}

.composer-shell {
  background: var(--glass-bg);
  backdrop-filter: blur(var(--glass-blur));
  -webkit-backdrop-filter: blur(var(--glass-blur));
  display: flex;
  align-items: center;
  gap: 10px;
  border-radius: 18px;
  border: 1px solid var(--glass-border);
  padding: 6px 8px 6px 6px;
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
}

.composer-shell:focus-within {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-glow);
}

/* Composer hover effect */
.composer-shell:hover {
  border-color: rgba(99, 102, 241, 0.5);
}

@media (prefers-reduced-motion: reduce) {
  .composer-shell:hover {
    border-color: inherit;
  }
}
```

- [ ] **Step 2: Test composer focus**

Click in message input
Expected: Glow border appears

- [ ] **Step 3: Test composer typing**

Type message, verify background
Expected: Glass effect maintained

- [ ] **Step 4: Commit**

```bash
git add frontend/css/style.css
git commit -m "style: apply glass effect to composer with hover effects"
```

---

### Task 11: Update Primary Button Styling

**Files:**
- Modify: `frontend/css/style.css:565-583` (btn.primary class)

- [ ] **Step 1: Update btn.primary class**

```css
.btn.primary {
  background: var(--primary);
  border-color: transparent;
  color: #ffffff;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.25);
  box-shadow: var(--glass-shadow-glow);
  transition: transform var(--transition-fast), background var(--transition-fast), box-shadow var(--transition-fast);
}

.btn.primary:hover:not(:disabled) {
  background: var(--primary-2);
  transform: scale(1.05);
  box-shadow: 0 0 0 4px var(--primary-glow);
}

.btn.primary:active:not(:disabled) {
  transform: scale(0.97);
}

html[data-theme="dark"] .btn.primary {
  background: var(--primary);
  color: #111827;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.25);
}

html[data-theme="dark"] .btn.primary:hover:not(:disabled) {
  background: var(--primary-2);
}

html[data-theme="dark"] .btn.primary:active:not(:disabled) {
  transform: scale(0.97);
}

@media (prefers-reduced-motion: reduce) {
  .btn.primary:hover:not(:disabled),
  .btn.primary:active:not(:disabled) {
    transform: none;
  }
}
```

- [ ] **Step 2: Test primary button rendering**

View primary buttons (new chat, send)
Expected: High contrast with text-shadow, scale on hover

- [ ] **Step 3: Test button hover**

Hover over primary button
Expected: Scale up to 1.05 with glow intensifies

- [ ] **Step 4: Test button active state**

Click and hold button
Expected: Scale down to 0.97

- [ ] **Step 5: Test reduced motion**

Enable reduced motion
Expected: No scale transform

- [ ] **Step 6: Commit**

```bash
git add frontend/css/style.css
git commit -m "style: update primary button with scale hover effect and text-shadow"
```

---

## Chunk 3: Background Animation System

### Task 12: Create Background Blob Markup

**Files:**
- Modify: `frontend/index.html:9-11` (after body tag)

- [ ] **Step 1: Add background blobs after body opening tag**

```html
<div class="background-blobs">
  <div class="background-blob-1" aria-hidden="true"></div>
  <div class="background-blob-2" aria-hidden="true"></div>
  <div class="background-blob-3" aria-hidden="true"></div>
</div>
```

- [ ] **Step 2: Test page loads without JS**

Open page with JavaScript disabled
Expected: Page renders, blobs visible but static

- [ ] **Step 3: Verify DOM structure**

Inspect page, check blobs container exists
Expected: background-blobs div with 3 blob children

- [ ] **Step 4: Commit**

```bash
git add frontend/index.html
git commit -m "feat: add background blob markup"
```

---

### Task 13: Create glass-ui.js Module

**Files:**
- Create: `frontend/js/glass-ui.js`

- [ ] **Step 1: Create glass-ui.js file**

```javascript
// glass-ui.js - Glass effects and animations

/**
 * Initialize background animation blobs
 */
export function initBackgroundAnimations() {
  // Animation is handled via CSS keyframes
  // This module can be extended for dynamic blob behavior
  console.log('Background animations initialized');
}

/**
 * Check if browser supports glass effects
 */
export function supportsBackdropFilter() {
  return CSS.supports('backdrop-filter', 'blur(20px)');
}

/**
 * Initialize glass UI enhancements
 */
export function init() {
  initBackgroundAnimations();

  // Log glass support
  if (!supportsBackdropFilter()) {
    console.warn('Backdrop filter not supported, using fallback styles');
  }
}
```

- [ ] **Step 2: Test module loads**

Open browser console
Expected: No errors, "Background animations initialized" logged

- [ ] **Step 3: Test supportsBackdropFilter**

In console, call glassUI.supportsBackdropFilter()
Expected: Returns true/false based on browser

- [ ] **Step 4: Commit**

```bash
git add frontend/js/glass-ui.js
git commit -m "feat: create glass-ui.js module"
```

---

### Task 14: Add Background Animation CSS

**Files:**
- Create: `frontend/css/components.css`

- [ ] **Step 1: Create components.css with blob animations**

```css
/* Background Animation Blobs */
.background-blobs {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
  pointer-events: none;
  overflow: hidden;
}

.background-blob-1,
.background-blob-2,
.background-blob-3 {
  position: fixed;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.8;
}

.background-blob-1 {
  width: 300px;
  height: 300px;
  background: rgba(99, 102, 241, 0.3);
  top: -10%;
  left: -10%;
  animation: float-blob-1 25s infinite ease-in-out;
}

.background-blob-2 {
  width: 400px;
  height: 400px;
  background: rgba(139, 92, 246, 0.25);
  top: 20%;
  right: -10%;
  animation: float-blob-2 30s infinite ease-in-out;
}

.background-blob-3 {
  width: 350px;
  height: 350px;
  background: rgba(34, 211, 238, 0.2);
  bottom: 10%;
  left: 50%;
  transform: translateX(-50%);
  animation: float-blob-3 22s infinite ease-in-out;
}

@keyframes float-blob-1 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(50px, 30px) scale(1.1); }
  50% { transform: translate(30px, 60px) scale(0.95); }
  75% { transform: translate(-20px, 40px) scale(1.05); }
}

@keyframes float-blob-2 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(-40px, 20px) scale(1.05); }
  50% { transform: translate(-20px, 40px) scale(0.9); }
  75% { transform: translate(-30px, 10px) scale(1.02); }
}

@keyframes float-blob-3 {
  0%, 100% { transform: translateX(-50%) translateY(0) scale(1); }
  33% { transform: translateX(-50%) translateY(30px) scale(1.1); }
  66% { transform: translateX(-50%) translateY(-20px) scale(0.95); }
}

@media (prefers-reduced-motion: reduce) {
  .background-blob-1,
  .background-blob-2,
  .background-blob-3 {
    animation: none;
  }
}
```

- [ ] **Step 2: Test blob animations**

Open app in browser
Expected: Blobs float smoothly across background

- [ ] **Step 3: Test reduced motion**

Enable reduced motion, reload
Expected: Blobs are static, no animation

- [ ] **Step 4: Commit**

```bash
git add frontend/css/components.css
git commit -m "feat: add background blob animations"
```

---

### Task 14b: Add components.css to HTML

**Files:**
- Modify: `frontend/index.html:7`

- [ ] **Step 1: Add components.css link**

```html
<link rel="stylesheet" href="./css/components.css" />
```

Place after style.css link.

- [ ] **Step 2: Test page loads**

Open app
Expected: Styles load, no errors in console

- [ ] **Step 3: Verify animations work**

Check that background blobs are animating
Expected: Smooth floating motion

- [ ] **Step 4: Commit**

```bash
git add frontend/index.html
git commit -m "feat: add components.css link to HTML"
```

---

## Chunk 4: Message Actions & Toast System

### Task 15: Create Toast Notification System

**Files:**
- Create: `frontend/js/features.js`

- [ ] **Step 1: Create features.js with toast system**

```javascript
// features.js - Toast notifications, export, message actions

const toastContainer = document.createElement('div');
toastContainer.id = 'toastContainer';
toastContainer.className = 'toast-container';
toastContainer.setAttribute('role', 'status');
toastContainer.setAttribute('aria-live', 'polite');
document.body.appendChild(toastContainer);

/**
 * Show toast notification
 */
export function showToast(message, type = 'info', duration = 4000) {
  const toast = document.createElement('div');
  toast.className = `toast toast-${type}`;
  toast.textContent = message;

  toastContainer.appendChild(toast);

  // Animate in
  requestAnimationFrame(() => {
    toast.classList.add('show');
  });

  // Auto dismiss
  setTimeout(() => {
    toast.classList.remove('show');
    setTimeout(() => toast.remove(), 200);
  }, duration);
}

/**
 * Initialize features module
 */
export function init() {
  console.log('Features module initialized');
}
```

- [ ] **Step 2: Create toast CSS**

Add to `frontend/css/components.css`:

```css
/* Toast Notifications */
.toast-container {
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 100;
  display: flex;
  flex-direction: column;
  gap: 10px;
  pointer-events: none;
}

.toast {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  opacity: 0;
  transform: translateY(20px);
  transition: opacity 200ms ease, transform 300ms cubic-bezier(0.16, 1, 0.3, 1);
  pointer-events: auto;
  min-width: 200px;
}

.toast.show {
  opacity: 1;
  transform: translateY(0);
}

.toast-success {
  border-left: 4px solid var(--success);
}

.toast-error {
  border-left: 4px solid var(--danger);
  background: rgba(239, 68, 68, 0.9);
  color: white;
}

.toast-info {
  border-left: 4px solid var(--primary);
}

@media (prefers-reduced-motion: reduce) {
  .toast {
    transition: opacity 0.01ms ease, transform 0.01ms cubic-bezier(0.16, 1, 0.3, 1);
  }
}
```

- [ ] **Step 3: Test toast notifications**

In browser console, run: `window.features.showToast('Test message', 'success')`
Expected: Toast appears, slides in from bottom

- [ ] **Step 4: Test toast types**

Run with 'error' and 'info' types
Expected: Different colors, correct styling

- [ ] **Step 5: Test auto-dismiss**

Wait 5 seconds
Expected: Toast fades out and is removed

- [ ] **Step 6: Commit**

```bash
git add frontend/css/components.css frontend/js/features.js
git commit -m "feat: add toast notification system"
```

---

### Task 16: Add Copy to Clipboard Function

**Files:**
- Modify: `frontend/js/features.js`

- [ ] **Step 1: Add copy function to features.js**

```javascript
/**
 * Copy text to clipboard
 */
export async function copyToClipboard(text, label = 'Content') {
  try {
    await navigator.clipboard.writeText(text);
    showToast(`${label} copied`, 'success');
    return true;
  } catch (err) {
    console.error('Copy failed:', err);
    showToast('Failed to copy', 'error');
    return false;
  }
}
```

- [ ] **Step 2: Test copy function**

In console: `window.features.copyToClipboard('Test message', 'Message')`
Expected: Success toast appears

- [ ] **Step 3: Verify clipboard content**

Paste from clipboard
Expected: "Test message" is pasted

- [ ] **Step 4: Test error handling**

Simulate clipboard error (use private window)
Expected: Error toast shown

- [ ] **Step 5: Commit**

```bash
git add frontend/js/features.js
git commit -m "feat: add copy to clipboard function"
```

---

### Task 17: Add Message Action Buttons

**Files:**
- Modify: `frontend/css/components.css`

- [ ] **Step 1: Add message action styles**

```css
/* Message Actions */
.msg {
  position: relative;
}

.msg-actions {
  position: absolute;
  top: -10px;
  right: 10px;
  display: flex;
  gap: 6px;
  opacity: 0;
  transform: translateY(-5px);
  transition: opacity 150ms ease, transform 150ms ease;
  pointer-events: none;
}

.msg:hover .msg-actions {
  opacity: 1;
  transform: translateY(0);
  pointer-events: auto;
}

.msg-action-btn {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  border: 1px solid var(--glass-border);
  background: var(--glass-bg);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  display: grid;
  place-items: center;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.msg-action-btn:hover {
  background: var(--primary);
  color: white;
  transform: scale(1.1);
}

.msg-action-btn:active {
  transform: scale(0.95);
}

.msg-action-btn svg {
  width: 14px;
  height: 14px;
}

@media (prefers-reduced-motion: reduce) {
  .msg-action-btn:hover,
  .msg-action-btn:active {
    transform: none;
  }
}
```

- [ ] **Step 2: Test action button hover**

Hover over a message bubble
Expected: Action buttons appear above bubble

- [ ] **Step 3: Test button clicks**

Click an action button
Expected: Visual feedback, scale on hover/active

- [ ] **Step 4: Commit**

```bash
git add frontend/css/components.css
git commit -m "style: add message action button styles"
```

---

### Task 18: Integrate Features Module

**Files:**
- Modify: `frontend/index.html:200-202` (before closing body)

- [ ] **Step 1: Add features.js script**

```html
<script type="module">
  import * as features from './js/features.js';
  import * as glassUI from './js/glass-ui.js';

  // Initialize modules
  window.features = features;
  window.glassUI = glassUI;

  features.init();
  glassUI.init();
</script>
```

Place before existing app.js script.

- [ ] **Step 2: Test modules load**

Open app, check console
Expected: "Features module initialized" and "Background animations initialized" logged

- [ ] **Step 3: Test toast from console**

Run: `window.features.showToast('Test', 'info')`
Expected: Toast appears

- [ ] **Step 4: Commit**

```bash
git add frontend/index.html
git commit -m "feat: integrate features and glass-ui modules"
```

---

## Chunk 5: Keyboard Shortcuts & Navigation

### Task 19: Create Keyboard Module

**Files:**
- Create: `frontend/js/keyboard.js`

- [ ] **Step 1: Create keyboard.js module**

```javascript
// keyboard.js - Keyboard shortcuts and navigation

/**
 * Initialize keyboard shortcuts
 */
export function init() {
  document.addEventListener('keydown', handleKeydown);
  console.log('Keyboard shortcuts initialized');
}

/**
 * Handle keyboard shortcuts
 */
function handleKeydown(e) {
  // Ctrl/Cmd + K: New chat
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault();
    window.newChat && window.newChat();
    return;
  }

  // Ctrl/Cmd + N: Create new chat
  if ((e.ctrlKey || e.metaKey) && e.key === 'n') {
    e.preventDefault();
    window.newChat && window.newChat();
    return;
  }

  // Ctrl/Cmd + /: Show shortcuts modal
  if ((e.ctrlKey || e.metaKey) && e.key === '/') {
    e.preventDefault();
    showShortcutsModal();
    return;
  }

  // Escape: Close modals/drawers, clear focus
  if (e.key === 'Escape') {
    closeAllModals();
    return;
  }
}

/**
 * Show shortcuts modal
 */
function showShortcutsModal() {
  showToast('Keyboard shortcuts: Ctrl+K (new chat), Ctrl+/ (help)', 'info');
}

/**
 * Close all modals and drawers
 */
function closeAllModals() {
  // Close auth dialog
  const authDialog = document.getElementById('authDialog');
  if (authDialog && authDialog.open) {
    authDialog.close();
  }

  // Close me dialog
  const meDialog = document.getElementById('meDialog');
  if (meDialog && meDialog.open) {
    meDialog.close();
  }

  // Clear focus from inputs
  document.activeElement && document.activeElement.blur();
}
```

- [ ] **Step 2: Test Ctrl+K shortcut**

Press Ctrl+K
Expected: New chat triggered (if available)

- [ ] **Step 3: Test Escape key**

Press Escape when dialog is open
Expected: Dialog closes

- [ ] **Step 4: Test Ctrl+/ shortcut**

Press Ctrl+/
Expected: Help toast appears

- [ ] **Step 5: Commit**

```bash
git add frontend/js/keyboard.js
git commit -m "feat: create keyboard shortcuts module"
```

---

### Task 19b: Integrate Keyboard Module

**Files:**
- Modify: `frontend/index.html:200-202`

- [ ] **Step 1: Add keyboard.js to import**

Update script block:

```html
<script type="module">
  import * as features from './js/features.js';
  import * as glassUI from './js/glass-ui.js';
  import * as keyboard from './js/keyboard.js';

  // Initialize modules
  window.features = features;
  window.glassUI = glassUI;
  window.keyboard = keyboard;

  features.init();
  glassUI.init();
  keyboard.init();
</script>
```

- [ ] **Step 2: Test all modules load**

Open app, check console
Expected: All three modules logged as initialized

- [ ] **Step 3: Test keyboard shortcuts**

Test Ctrl+K, Ctrl+/, Escape
Expected: All shortcuts work

- [ ] **Step 4: Commit**

```bash
git add frontend/index.html
git commit -m "feat: integrate keyboard module"
```

---

## Chunk 6: Mobile Responsive Design

### Task 20: Add Mobile Breakpoint Styles

**Files:**
- Modify: `frontend/css/style.css` (append at end)

- [ ] **Step 1: Add mobile breakpoint media query**

```css
/* Mobile Responsive (< 768px) */
@media (max-width: 767px) {
  .app {
    grid-template-columns: 1fr;
  }

  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    width: 280px;
    z-index: 50;
    transform: translateX(-100%);
    transition: transform 300ms cubic-bezier(0.16, 1, 0.3, 1);
  }

  .sidebar.open {
    transform: translateX(0);
  }

  .sidebar-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(4px);
    -webkit-backdrop-filter: blur(4px);
    z-index: 45;
    opacity: 0;
    pointer-events: none;
    transition: opacity 200ms ease;
  }

  .sidebar-backdrop.open {
    opacity: 1;
    pointer-events: auto;
  }

  /* Touch targets */
  .icon-btn {
    min-width: 44px;
    min-height: 44px;
    padding: 12px;
  }

  /* Composer fixed on mobile */
  .composer {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    padding-bottom: calc(14px + env(safe-area-inset-bottom));
    z-index: 30;
  }

  .composer-shell {
    margin: 0 16px;
    border-radius: 18px;
  }

  /* Top bar compact on mobile */
  .topbar {
    height: 56px;
    padding-top: calc(16px + env(safe-area-inset-top));
  }

  /* Messages full height on mobile */
  .messages {
    padding-bottom: 140px; /* Space for fixed composer */
  }
}
```

- [ ] **Step 2: Test on mobile viewport**

Resize browser to < 768px
Expected: Sidebar hidden, composer fixed at bottom

- [ ] **Step 3: Test sidebar backdrop**

Add sidebar-backdrop div to HTML, test toggle
Expected: Backdrop fades in/out

- [ ] **Step 4: Commit**

```bash
git add frontend/css/style.css
git commit -m "style: add mobile responsive breakpoint styles"
```

---

### Task 21: Add Hamburger Menu Button

**Files:**
- Modify: `frontend/index.html:40-48` (topbar section)

- [ ] **Step 1: Add hamburger button to topbar**

```html
<button id="menuBtn" class="icon-btn menu-btn" type="button" aria-label="Open menu">
  <svg class="menu-icon" viewBox="0 0 24 24" aria-hidden="true">
    <path d="M3 12h18M3 6h18M3 18h18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
  </svg>
</button>
```

Place at beginning of topbar content.

- [ ] **Step 2: Add sidebar backdrop**

```html
<div id="sidebarBackdrop" class="sidebar-backdrop" aria-hidden="true"></div>
```

Place before app closing tag.

- [ ] **Step 3: Test hamburger button visibility**

Resize to mobile, verify button appears
Expected: Hamburger icon visible in top-left

- [ ] **Step 4: Commit**

```bash
git add frontend/index.html
git commit -m "feat: add hamburger menu button and sidebar backdrop"
```

---

### Task 22: Add Mobile Menu Toggle Functionality

**Files:**
- Modify: `frontend/js/keyboard.js`

- [ ] **Step 1: Add mobile menu functions to keyboard.js**

```javascript
/**
 * Toggle mobile sidebar
 */
export function toggleSidebar() {
  const sidebar = document.querySelector('.sidebar');
  const backdrop = document.getElementById('sidebarBackdrop');

  if (sidebar && backdrop) {
    sidebar.classList.toggle('open');
    backdrop.classList.toggle('open');

    // Trap focus in sidebar when open
    if (sidebar.classList.contains('open')) {
      const focusable = sidebar.querySelector('button, input, [tabindex]');
      if (focusable) focusable.focus();
    }
  }
}

/**
 * Close mobile sidebar
 */
export function closeSidebar() {
  const sidebar = document.querySelector('.sidebar');
  const backdrop = document.getElementById('sidebarBackdrop');

  if (sidebar && backdrop) {
    sidebar.classList.remove('open');
    backdrop.classList.remove('open');
  }
}

// Update init to add menu button listener
const originalInit = init;
init = function() {
  originalInit();

  const menuBtn = document.getElementById('menuBtn');
  if (menuBtn) {
    menuBtn.addEventListener('click', toggleSidebar);
  }

  const backdrop = document.getElementById('sidebarBackdrop');
  if (backdrop) {
    backdrop.addEventListener('click', closeSidebar);
  }
};
```

- [ ] **Step 2: Test sidebar toggle**

Resize to mobile, click hamburger button
Expected: Sidebar slides in from left

- [ ] **Step 3: Test backdrop click**

Click backdrop
Expected: Sidebar slides out

- [ ] **Step 4: Test Escape to close**

Open sidebar, press Escape
Expected: Sidebar closes

- [ ] **Step 5: Commit**

```bash
git add frontend/js/keyboard.js
git commit -m "feat: add mobile sidebar toggle functionality"
```

---

### Task 23: Add Mobile Menu Styles

**Files:**
- Modify: `frontend/css/components.css`

- [ ] **Step 1: Add menu button styles**

```css
/* Mobile Menu Button */
.menu-btn {
  display: none;
}

.menu-icon {
  width: 20px;
  height: 20px;
  stroke: var(--text-primary);
}

@media (max-width: 767px) {
  .menu-btn {
    display: grid;
    place-items: center;
  }
}
```

- [ ] **Step 2: Test menu button on mobile**

Resize to mobile, check hamburger button
Expected: Visible and properly styled

- [ ] **Step 3: Test menu button hover**

Hover over hamburger button
Expected: Subtle scale effect

- [ ] **Step 4: Commit**

```bash
git add frontend/css/components.css
git commit -m "style: add mobile menu button styles"
```

---

## Chunk 7: Polish & Testing

### Task 24: Update App Elements Reference

**Files:**
- Modify: `frontend/js/app.js:1-30` (els object)

- [ ] **Step 1: Add new elements to els object**

```javascript
const els = {
  sessionsList: document.getElementById("sessionsList"),
  modelSelect: document.getElementById("modelSelect"),
  newChatBtn: document.getElementById("newChatBtn"),
  themeToggle: document.getElementById("themeToggle"),
  meBtn: document.getElementById("meBtn"),
  meDialog: document.getElementById("meDialog"),
  meClose: document.getElementById("meClose"),
  meAvatar: document.getElementById("meAvatar"),
  meName: document.getElementById("meName"),
  meAvatar2: document.getElementById("meAvatar2"),
  meName2: document.getElementById("meName2"),
  meId: document.getElementById("meId"),
  authDialog: document.getElementById("authDialog"),
  authClose: document.getElementById("authClose"),
  loginForm: document.getElementById("loginForm"),
  registerForm: document.getElementById("registerForm"),
  loginEmail: document.getElementById("loginEmail"),
  loginPassword: document.getElementById("loginPassword"),
  registerName: document.getElementById("registerName"),
  registerEmail: document.getElementById("registerEmail"),
  registerPassword: document.getElementById("registerPassword"),
  messages: document.getElementById("messages"),
  prompt: document.getElementById("prompt"),
  sendBtn: document.getElementById("sendBtn"),
  errorBar: document.getElementById("errorBar"),
  fileInput: document.getElementById("fileInput"),
  attachments: document.getElementById("attachments"),
  logoutBtn: document.getElementById("logoutBtn"),
  // New elements
  menuBtn: document.getElementById("menuBtn"),
  sidebarBackdrop: document.getElementById("sidebarBackdrop"),
};
```

- [ ] **Step 2: Test element references**

Open app, check no console errors
Expected: All elements found, no null reference errors

- [ ] **Step 3: Commit**

```bash
git add frontend/js/app.js
git commit -m "refactor: update app elements reference with mobile elements"
```

---

### Task 25: Add Message Entry Animation

**Files:**
- Modify: `frontend/css/components.css`

- [ ] **Step 1: Add message entry animation**

```css
/* Message Entry Animation */
.msg {
  animation: messageSlideIn 300ms cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes messageSlideIn {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@media (prefers-reduced-motion: reduce) {
  .msg {
    animation: none;
    opacity: 1;
    transform: none;
  }
}
```

- [ ] **Step 2: Test message animation**

Send a message
Expected: Message slides up with fade effect

- [ ] **Step 3: Test reduced motion**

Enable reduced motion, send message
Expected: No animation, instant appearance

- [ ] **Step 4: Commit**

```bash
git add frontend/css/components.css
git commit -m "style: add message entry animation"
```

---

### Task 26: Add Typing Indicator

**Files:**
- Modify: `frontend/css/components.css`

- [ ] **Step 1: Add typing indicator styles**

```css
/* Typing Indicator */
.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 8px 0;
}

.typing-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--text-secondary);
  animation: typingBounce 1.4s infinite ease-in-out;
}

.typing-dot:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-dot:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typingBounce {
  0%, 60%, 100% {
    transform: translateY(0);
  }
  30% {
    transform: translateY(-8px);
  }
}

@media (prefers-reduced-motion: reduce) {
  .typing-dot {
    animation: none;
  }
}
```

- [ ] **Step 2: Test typing indicator**

Add typing-indicator HTML to a message
Expected: Dots bounce smoothly

- [ ] **Step 3: Test reduced motion**

Enable reduced motion
Expected: Dots are static

- [ ] **Step 4: Commit**

```bash
git add frontend/css/components.css
git commit -m "style: add typing indicator animation"
```

---

### Task 27: Add Hover State Animations

**Files:**
- Modify: `frontend/css/style.css`

- [ ] **Step 1: Update button hover animations**

Find and replace .icon-btn:hover styles:

```css
.icon-btn:hover:not(:disabled) {
  background: var(--bg-hover);
  border-color: var(--border-hover);
  color: var(--text-primary);
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.2);
}

.icon-btn:active:not(:disabled) {
  transform: scale(0.97);
}

@media (prefers-reduced-motion: reduce) {
  .icon-btn:hover:not(:disabled),
  .icon-btn:active:not(:disabled) {
    transform: none;
  }
}
```

- [ ] **Step 2: Test button hover**

Hover over various buttons
Expected: Smooth scale animation with glow

- [ ] **Step 3: Test button active state**

Click and hold button
Expected: Scale down to 0.97

- [ ] **Step 4: Commit**

```bash
git add frontend/css/style.css
git commit -m "style: enhance button hover/active animations"
```

---

### Task 28: Final Cross-Browser Testing

**Files:**
- Test: All changes across browsers

- [ ] **Step 1: Test in Chrome (90+)**

Open app in Chrome
Expected: All features work, glass effects visible

- [ ] **Step 2: Test in Firefox (103+)**

Open app in Firefox
Expected: All features work (blur effects require FF 103+)

- [ ] **Step 3: Test in Safari (14+)**

Open app in Safari
Expected: All features work

- [ ] **Step 4: Test in Edge (90+)**

Open app in Edge
Expected: All features work

- [ ] **Step 5: Document any browser-specific issues**

Create notes if any browser has issues
Expected: Fallbacks work for older browsers

---

### Task 29: Accessibility Audit

**Files:**
- Test: Accessibility features

- [ ] **Step 1: Test keyboard navigation**

Tab through entire interface
Expected: Logical tab order, all interactive elements reachable

- [ ] **Step 2: Test focus visibility**

Verify all focus states are visible
Expected: Clear focus rings on all interactive elements

- [ ] **Step 3: Test screen reader compatibility**

Use screen reader (NVDA, VoiceOver, JAWS)
Expected: All interactive elements announced correctly

- [ ] **Step 4: Test contrast ratios**

Use contrast checker tool on both themes
Expected: All text meets WCAG AA (4.5:1)

- [ ] **Step 5: Test reduced motion**

Enable OS reduced motion preference
Expected: All animations disabled, instant transitions

- [ ] **Step 6: Test high contrast mode**

Enable OS high contrast preference
Expected: High opacity backgrounds, solid borders

- [ ] **Step 7: Commit any accessibility fixes**

```bash
git add frontend/css/style.css frontend/css/components.css
git commit -m "fix: accessibility improvements based on audit"
```

---

### Task 30: Performance Optimization

**Files:**
- Test: Animation performance

- [ ] **Step 1: Check animation performance**

Open Chrome DevTools Performance tab
Expected: Animations run at 60fps

- [ ] **Step 2: Optimize if needed**

If animations drop below 60fps:
- Reduce blob animation complexity
- Simplify hover effects
Expected: Maintained 60fps performance

- [ ] **Step 3: Test on lower-end devices**

Test on slower device
Expected: Smooth performance, acceptable load times

- [ ] **Step 4: Commit any performance fixes**

```bash
git add frontend/css/style.css frontend/css/components.css
git commit -m "perf: optimize animation performance"
```

---

## Completion

All tasks completed. The Glass & Futuristic Luban UI transformation is ready for deployment.

### Summary of Changes

- CSS variables updated with glass system including bubble colors
- Glass panel effects applied to all major components (sidebar, session cards, messages, top bar, composer)
- Hover/active animations implemented for all interactive elements
- Animated gradient background with 3 floating blobs
- Toast notification system implemented
- Keyboard shortcuts module created (Ctrl+K, Ctrl+N, Ctrl+/, Escape)
- Mobile responsive design with sidebar drawer and hamburger menu
- Accessibility improvements (focus, reduced motion, high contrast)
- Cross-browser compatibility with fallbacks
- Message entry animation and typing indicator

### Final Verification

- [ ] Visual impact: Glass effects are visible and modern
- [ ] Readability: All text meets WCAG AA contrast (4.5:1)
- [ ] Performance: Animations run at 60fps
- [ ] Accessibility: Keyboard navigation, screen reader support, reduced motion, high contrast
- [ ] Responsive: Mobile, tablet, desktop layouts work
- [ ] Browser support: Works in Chrome, Firefox, Safari, Edge
- [ ] Code quality: Files are focused, responsibilities clear, no duplication

---

**Plan complete and saved to `docs/superpowers/plans/2026-03-17-glass-futuristic-ui-plan.md`. Ready to execute?**
