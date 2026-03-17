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
| `--glass-bg` | `rgba(255, 255, 255, 0.25)` | `rgba(30, 41, 59, 0.4)` | Frosted glass panels |
| `--glass-border` | `rgba(255, 255, 255, 0.4)` | `rgba(255, 255, 255, 0.1)` | Glass edge highlight |
| `--text-primary` | `#1e293b` | `#f1f5f9` | Primary text |
| `--text-secondary` | `#64748b` | `#94a3b8` | Secondary text |
| `--primary` | `#6366f1` | `#818cf8` | Brand indigo accent |
| `--primary-glow` | `rgba(99, 102, 241, 0.4)` | `rgba(129, 140, 248, 0.5)` | Glow effect |
| `--bubble-user` | `linear-gradient(135deg, #6366f1, #8b5cf6)` | `linear-gradient(135deg, #4f46e5, #7c3aed)` | User message gradient |
| `--bubble-assistant` | `rgba(255, 255, 255, 0.5)` | `rgba(51, 65, 85, 0.5)` | Assistant glass |

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
- **Composer:** Floating glass panel (not attached to bottom) with blur

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
Level 30: Composer floating panel
Level 40: Modals / Dialogs
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
- **Animated Gradient Blobs:** 3-4 colored orbs floating slowly
- **Duration:** 20-30s per loop
- **Motion:** Slow translate + scale oscillation
- **Respects Reduced Motion:** Pauses or becomes static

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
}
```

---

## 4. Feature Additions & Enhancements

### New Components

#### 1. Message Actions (Hover on bubbles)
- **Copy Button:** Copy message to clipboard
- **Regenerate Button:** Re-generate assistant response
- **Edit Button:** Edit user message and retry
- **Delete Button:** Delete message from conversation
- *Position:* Top-right of each bubble, visible on hover

#### 2. Code Syntax Highlighting
- **Code Blocks:** Glass panel with darker background
- **Syntax Highlighting:** Highlight.js or Prism.js
- **Language Detection:** Auto-detect from markdown fences
- **Copy Button:** Dedicated copy for code blocks
- **Line Numbers:** Optional toggle in settings

#### 3. Model Information Card
- **Hover on Model Select:** Shows model details
  - Model name and description
  - Token limit
  - Response speed indicator
  - Capabilities list

#### 4. Attachment Preview
- **Thumbnail:** Glass-framed preview for images/videos
- **File Type Icon:** For non-visual files
- **Remove Button:** X icon on attachment
- **Size Badge:** File size display

#### 5. Export Options
- **Export Button:** In top bar
- **Formats:** Markdown, PDF, JSON
- **Scope:** Current session or all sessions

#### 6. Search Sessions
- **Search Input:** Above session list in sidebar
- **Live Filtering:** Real-time search as you type
- **Highlight:** Matching text in results

#### 7. Quick Actions (Keyboard Shortcuts)
- **Ctrl/Cmd + K:** Start new chat
- **Ctrl/Cmd + /:** Show shortcuts modal
- **Ctrl/Cmd + N:** Create new chat
- **Escape:** Close composer focus / close modals

#### 8. Toast Notifications
- **Position:** Bottom-right corner
- **Types:** Success (green), Error (red), Info (blue)
- **Auto-dismiss:** 4 seconds
- **Manual dismiss:** Click to close

#### 9. Typing Indicator
- **Appearance:** Glass capsule with 3 animated dots
- **Position:** Below last user message
- **Animation:** Bouncing dots with 150ms stagger

#### 10. Welcome/Onboarding Experience
- **Empty State:** Animated Luban mascot illustration
- **Suggested Prompts:** Glass cards with quick-start questions
- **Feature Tour:** Optional guided tour on first visit

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

### Mobile-Specific Adaptations

#### 1. Sidebar as Slide-Over Drawer
- **Trigger:** Hamburger menu (≡) in top-left
- **Behavior:** Full-height glass panel slides from left
- **Backdrop:** Dimmed glass overlay
- **Close:** Tap backdrop, close button, or swipe right

#### 2. Touch Targets
- **Minimum size:** 44×44px (iOS standard)
- **Button padding:** Increased to 12px
- **Tap area:** Extended beyond visual bounds using hit regions

#### 3. Composer on Mobile
- **Fixed position:** Always visible at bottom
- **Auto-resize:** Grows to max 120px, then scrolls
- **Keyboard handling:** Adjusts for virtual keyboard
- **Send button:** Prominent right-side placement

#### 4. Message Actions
- **Long press:** Opens action menu (copy, edit, delete)
- **Swipe left:** Quick delete option
- **Menu:** Bottom sheet glass panel

#### 5. Model Selector
- **Position:** Compact dropdown in top bar
- **Tap to expand:** Full-width glass sheet with model details

#### 6. Safe Area Support
- **Top bar:** Respects status bar notch
- **Bottom composer:** Respects home indicator
- **Side padding:** Adjusts for device edges

#### 7. Scroll Behavior
- **Overscroll:** Bounce at edges (iOS style)
- **Pull to refresh:** Not needed (real-time updates)
- **Scroll momentum:** Native smooth scrolling

#### 8. Gesture Support
- **Swipe left from edge:** Open sidebar
- **Swipe down on drawer:** Close drawer
- **Two-finger tap:** Show message actions

### Typography Scaling

| Screen | Body Text | Heading |
|--------|-----------|---------|
| Mobile | 14px | 16px |
| Tablet | 15px | 18px |
| Desktop | 16px | 20px |

### Orientation Support
- **Landscape:** Full-width messages, compact top bar
- **Portrait:** Standard layout

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
}
```

### Glass Panel Mixin
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

### Key Components Structure

```
frontend/
├── css/
│   ├── style.css           # Main stylesheet (update with glass effects)
│   └── components.css      # New: Glass component styles
├── js/
│   ├── app.js              # Main app logic
│   ├── glass-ui.js         # New: Glass effects and animations
│   └── features.js         # New: Enhanced features (copy, export, etc.)
└── index.html             # Update with new markup
```

### External Dependencies
- **Highlight.js or Prism.js:** Code syntax highlighting
- **Inter Font:** Google Fonts (already included)
- **Lucide Icons or Heroicons:** SVG icons (already used)

---

## 7. Accessibility Checklist

- [ ] Minimum 4.5:1 contrast ratio for text (WCAG AA)
- [ ] Visible focus rings on all interactive elements
- [ ] Keyboard navigation support for all features
- [ ] Aria-labels on icon-only buttons
- [ ] Reduced motion support for animations
- [ ] Touch targets minimum 44×44px on mobile
- [ ] Screen reader friendly markup
- [ ] High contrast mode support
- [ ] Skip to main content link
- [ ] Form error messages associated with inputs

---

## 8. Implementation Priority

### Phase 1: Core Glass Transformation
1. Update CSS with glass variables and effects
2. Apply glass styling to existing components
3. Implement animated gradient background
4. Update typography to Inter font
5. Add hover/active state animations

### Phase 2: Layout Refinements
1. Refine sidebar glass panel styling
2. Update composer as floating glass panel
3. Add message action buttons
4. Implement model information card
5. Add typing indicator

### Phase 3: Feature Additions
1. Implement code syntax highlighting
2. Add copy/edit/delete message actions
3. Build export functionality
4. Add session search
5. Implement keyboard shortcuts
6. Add toast notifications

### Phase 4: Mobile Experience
1. Implement sidebar drawer
2. Add touch-specific interactions
3. Optimize composer for mobile
4. Add gesture support
5. Safe area handling

---

## 9. Success Criteria

The Glass & Futuristic redesign is successful when:

1. **Visual Impact:** The interface feels distinctly modern with glass effects and animated backgrounds
2. **Readability:** All text meets WCAG AA contrast standards in both light and dark modes
3. **Performance:** Animations run at 60fps with no jank
4. **Accessibility:** All features are keyboard accessible and screen reader friendly
5. **Responsive:** Experience is polished on mobile, tablet, and desktop
6. **Feature Parity:** All existing features work with the new design
7. **Browser Support:** Works in modern browsers (Chrome, Firefox, Safari, Edge)

---

**Document Status:** Ready for implementation
**Next Step:** Create detailed implementation plan
