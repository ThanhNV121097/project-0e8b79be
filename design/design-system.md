# Design System — Todo List App v2

> Source of truth: the approved `index.html`.
> Every value below is extracted from it. Changing a value here without
> changing the approved design is a defect.

Last updated: 2026-08-11

## 1. Foundations

### 1.1 Color

Semantic tokens. Name by job, never by hue.

| Token | Value | Used for |
|---|---|---|
| `--color-bg` | `#F8FAFC` | Page background, input background, empty-state panels, mini-state surfaces |
| `--color-surface` | `#FFFFFF` | Cards, navigation menu, links, todo rows, buttons without fill |
| `--color-surface-raised` | `#FFFFFF` | Mobile navigation popover and raised surfaces |
| `--color-text` | `#0F172A` | Body text and default UI text |
| `--color-heading` | `#0B1220` | Hero h1 text |
| `--color-text-muted` | `#64748B` | Secondary text, captions, inactive tabs and links |
| `--color-placeholder` | `#94A3B8` | Input placeholder and completed todo text |
| `--color-primary` | `#2563EB` | Primary actions, brand emphasis, active controls, focus source color |
| `--color-primary-hover` | `#1D4ED8` | Primary button hover background |
| `--color-primary-soft` | `#EFF6FF` | Primary hover wash, active tab background, eyebrow and icon backgrounds |
| `--color-primary-tint` | `#DBEAFE` | Navigation border, tab border, page radial background start |
| `--color-primary-border` | `#BFDBFE` | Input border, dashed empty border |
| `--color-primary-control` | `#93C5FD` | Checkbox border and loading dots |
| `--color-primary-gradient-end` | `#60A5FA` | Logo gradient end |
| `--color-primary-text` | `#FFFFFF` | Text on primary and destructive action backgrounds |
| `--color-success` | `#10B981` | Completed checkbox background and border |
| `--color-success-text` | `#065F46` | Saved-state pill text |
| `--color-success-bg` | `#D1FAE5` | Saved-state pill background |
| `--color-warning` | `#F59E0B` | Warning token present in root for warning state |
| `--color-warning-bg` | `#FEF3C7` | Sync notice background |
| `--color-warning-border` | `#FDE68A` | Sync notice border |
| `--color-warning-text` | `#92400E` | Sync notice text |
| `--color-danger` | `#EF4444` | Delete action, validation and error text |
| `--color-danger-bg` | `#FEF2F2` | Delete button and error mini-state background |
| `--color-danger-hover-bg` | `#FEE2E2` | Delete button hover background |
| `--color-border` | `#E2E8F0` | Card, todo row, state card, panel borders |
| `--color-border-soft` | `#E0ECFF` | Stat card border |
| `--color-focus` | `rgba(37,99,235,.35)` | Visible keyboard focus ring |
| `--color-nav-bg` | `rgba(248,250,252,.82)` | Sticky navigation background |
| `--color-nav-border` | `rgba(219,234,254,.9)` | Sticky navigation border |
| `--color-card-bg` | `rgba(255,255,255,.92)` | Main todo app card background |
| `--color-stat-bg` | `rgba(255,255,255,.72)` | Stat card translucent background |
| `--color-primary-shadow` | `rgba(37,99,235,.14)` | Main card shadow |
| `--color-primary-shadow-strong` | `rgba(37,99,235,.24)` | Logo and primary button shadow |
| `--color-neutral-shadow` | `rgba(15,23,42,.06)` | State card shadow |
| `--color-neutral-shadow-strong` | `rgba(15,23,42,.08)` | Mobile toggle shadow |

#### Contrast audit

Every text-on-background pair actually used. Body text ≥ 4.5:1, large text (≥ 18.66px bold or ≥ 24px) ≥ 3:1, UI borders ≥ 3:1.

| Foreground | Background | Ratio | Passes |
|---|---|---|---|
| `--color-text` `#0F172A` | `--color-bg` `#F8FAFC` | `16.8:1` | AA |
| `--color-heading` `#0B1220` | `--color-bg` `#F8FAFC` | `17.5:1` | AA |
| `--color-text-muted` `#64748B` | `--color-bg` `#F8FAFC` | `4.8:1` | AA |
| `--color-text-muted` `#64748B` | `--color-surface` `#FFFFFF` | `4.8:1` | AA |
| `--color-placeholder` `#94A3B8` | `--color-bg` `#F8FAFC` | `2.5:1` | FAIL for body text; acceptable only as placeholder support text |
| `--color-primary` `#2563EB` | `--color-primary-soft` `#EFF6FF` | `5.2:1` | AA |
| `--color-primary` `#2563EB` | `--color-surface` `#FFFFFF` | `5.2:1` | AA |
| `--color-primary-text` `#FFFFFF` | `--color-primary` `#2563EB` | `5.2:1` | AA |
| `--color-primary-text` `#FFFFFF` | `--color-primary-hover` `#1D4ED8` | `6.7:1` | AA |
| `--color-primary-text` `#FFFFFF` | `--color-danger` `#EF4444` | `3.8:1` | AA Large only; FAIL for normal body text |
| `--color-success-text` `#065F46` | `--color-success-bg` `#D1FAE5` | `7.4:1` | AA |
| `--color-warning-text` `#92400E` | `--color-warning-bg` `#FEF3C7` | `5.9:1` | AA |
| `--color-danger` `#EF4444` | `--color-danger-bg` `#FEF2F2` | `3.5:1` | AA Large only; FAIL for normal body text |
| `--color-placeholder` `#94A3B8` | `--color-surface` `#FFFFFF` | `2.5:1` | FAIL for completed todo text if read as body content |
| `--color-primary-control` `#93C5FD` | `--color-surface` `#FFFFFF` | `1.8:1` | FAIL as UI border contrast |
### 1.2 Spacing

Base unit: `2px`. The approved design mostly uses a 2px-derived scale with several component-specific values.

| Token | Value |
|---|---|
| `--space-1` | `4px` |
| `--space-2` | `8px` |
| `--space-3` | `10px` |
| `--space-4` | `12px` |
| `--space-5` | `14px` |
| `--space-6` | `16px` |
| `--space-7` | `18px` |
| `--space-8` | `20px` |
| `--space-9` | `22px` |
| `--space-10` | `24px` |
| `--space-11` | `28px` |
| `--space-12` | `30px` |
| `--space-13` | `34px` |
| `--space-14` | `46px` |
| `--space-15` | `62px` |
| `--space-16` | `70px` |
| `--space-17` | `80px` |
| `--space-hairline` | `1px` |
| `--space-sr` | `1px` |

Component-specific spacing pairs from the approved CSS:

| Use | Value |
|---|---|
| Navigation inner padding | `16px 22px` |
| Main padding | `34px 22px 80px` |
| Hero padding | `46px 0 28px` |
| App card padding | `22px` |
| Primary button padding | `12px 16px` |
| Ghost link padding | `10px 12px` |
| Input padding | `15px 16px` |
| Todo item padding | `12px` |
| State card padding | `20px` |
| Accordion padding | `16px` |
| Panel padding | `0 16px` |
| Empty state padding | `30px 12px` |

### 1.3 Typography

Font families (include the fallback stack and how the font is loaded):

- Body: `Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif`; Inter is referenced in CSS and falls back to system UI if not installed.
- Headings: inherit the body stack.
- Mono: no mono family is defined or used in the approved design.

| Token | Size | Line height | Weight | Used for |
|---|---|---|---|---|
| `--text-xs` | `12px` | normal | `800` | Saved-state pill |
| `--text-sm` | `13px` | normal | regular | Stat captions, validation error |
| `--text-base` | browser default `16px` | normal / `1.65` in paragraphs | regular | Body, controls, todo text |
| `--text-summary` | `14px` | normal | regular | Todo summary row |
| `--text-lead` | `19px` | `1.65` | regular | Hero paragraph |
| `--text-stat` | `24px` | normal | bold | Stat numbers |
| `--text-app-title` | `28px` | normal | default heading weight | Todo card h2 |
| `--text-section-title` | `clamp(30px,4vw,46px)` | normal | default heading weight | Section h2 |
| `--text-hero-title` | `clamp(38px,6vw,68px)` | `.96` | default heading weight | Hero h1 |
| `--text-mini-icon` | `28px`, `30px`, `34px` | normal | `900` for error mini | Empty icon, error mini icon, empty mini icon |

Heading levels are used in order and never skipped for visual sizing: h1 for the page hero, h2 for app card and sections, h3 for empty and state cards.

### 1.4 Radius, border, shadow, motion

| Token | Value | Used for |
|---|---|---|
| `--radius-xs` | `10px` | Skip link |
| `--radius-sm` | `12px` | Logo, nav links, mobile toggle, destructive icon button, danger button |
| `--radius-md` | `14px` | Primary button, sync notice |
| `--radius-lg` | `16px` | Input, todo row, panel, accordion |
| `--radius-xl` | `18px` | Stat card, mobile nav, empty state |
| `--radius-2xl` | `20px` | Empty icon |
| `--radius-card` | `22px` | State card |
| `--radius-app` | `24px` | Main app card |
| `--radius-full` | `999px` / `9999px` | Eyebrow, pill, tabs, checkbox, loading dots |
| `--border-width` | `1px` | Default card, input, divider, tab, panel and notice borders |
| `--border-width-strong` | `2px` | Todo completion checkbox border |
| `--shadow-sm` | `0 6px 18px rgba(15,23,42,.08)` | Mobile menu toggle |
| `--shadow-md` | `0 10px 24px rgba(37,99,235,.24)` | Logo mark |
| `--shadow-lg` | `0 14px 28px rgba(37,99,235,.24)` | Primary button |
| `--shadow-xl` | `0 14px 38px rgba(15,23,42,.06)` | State cards |
| `--shadow-app` | `0 24px 70px rgba(37,99,235,.14)` | Main app card and mobile nav |
| `--duration-fast` | `.18s` | Primary button hover transform and background |
| `--duration-task` | `.22s` | Todo item insertion animation |
| `--duration-card` | `.6s` | App card entrance animation |
| `--duration-loading` | `.8s` | Loading dot pulse animation |
| `--easing` | `ease` | All declared keyframe animations and transitions |

Motion does not currently include a `prefers-reduced-motion: reduce` override in the approved HTML.

### 1.5 Layout and breakpoints

| Name | Min width | Container | Columns | Gutter |
|---|---|---|---|---|
| `base` | `0px` | Fluid width with `22px` page side padding | Single column where responsive media query applies | `10px` to `28px` depending component |
| `md` | `820px` max-width media query boundary | `1120px` max container | Hero and details use two columns; state grid uses three columns | `18px`, `28px`, `34px` |
| `lg` | Not separately defined | `1120px` max container | Same as `md` above boundary | Same as above |
| `xl` | Not separately defined | `1120px` max container | Same as `md` above boundary | Same as above |

Z-index scale (only these values are allowed):

| Layer | Value |
|---|---|
| Base | `0` |
| Sticky header | `5` |
| Dropdown | `5` within the sticky header stacking context |
| Modal backdrop | Not used in approved design |
| Modal | Not used in approved design |
| Toast | Not used in approved design |
| Skip link focus | `10` |
## 2. Components

One subsection per reusable component. Every component lists all states.

### 2.1 Primary Button / Primary Link

**Purpose** — Use for the main action on a screen or section, such as opening the app or adding a task. Do not use for secondary navigation or destructive actions.

**Anatomy** — `[label]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Primary fill | `--color-primary`, `--color-primary-text`, `--radius-md`, `--shadow-lg` | Main action buttons and primary anchor links |
| Danger fill | `--color-danger`, `--color-primary-text`, `--radius-sm` | Destructive confirmation button pattern if needed |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Default | content height from `12px 16px` padding | `12px 16px` | `--text-base` |
| Compact danger | content height from `10px 12px` padding | `10px 12px` | `--text-base` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Blue fill, white text, soft blue shadow | `--color-primary`, `--color-primary-text`, `--shadow-lg` |
| Hover | Darker blue fill and `translateY(-1px)` lift | `--color-primary-hover`, `--duration-fast` |
| Focus (keyboard) | Visible 3px blue translucent focus ring with 3px offset | `--color-focus` |
| Active / pressed | No separate active style is defined; remains at hover/default visual treatment | `--color-primary` |
| Disabled | No disabled style is defined; if disabled is required, keep the component non-interactive and do not invent new colors outside this system | Existing primary tokens only |
| Loading | No loading style is defined for this component; keep label stable and pair with the loading-state component when needed | Existing primary tokens only |
| Error | Use danger fill variant only for destructive/error action confirmation | `--color-danger`, `--color-primary-text` |
| Empty | In empty states, primary action copy should explain the next action; the button keeps primary fill | `--color-primary`, `--color-primary-text` |

**Accessibility** — Native `button` or anchor when navigating. Keyboard focus must show the focus ring. Minimum hit target is at least 44×44px through padding and line height.

### 2.2 Ghost Navigation Link / Secondary Button

**Purpose** — Use for secondary navigation, non-primary actions, and low-emphasis controls such as simulating an error.

**Anatomy** — `[label]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Ghost | `--color-text-muted`, transparent background, `--radius-sm` | Secondary links and buttons |
| Hover wash | `--color-primary-soft`, `--color-primary` | Hover indication |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Default | content height from `10px 12px` padding | `10px 12px` | `--text-base` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Transparent background, muted text | `--color-text-muted` |
| Hover | Pale primary background and primary text | `--color-primary-soft`, `--color-primary` |
| Focus (keyboard) | Visible 3px blue translucent focus ring with 3px offset | `--color-focus` |
| Active / pressed | No separate active style is defined; remains hover/default | Existing ghost tokens |
| Disabled | No disabled style is defined; non-interactive ghost controls should keep muted text without hover behavior | `--color-text-muted` |
| Loading | No loading style is defined; pair with loading-state component if async work is running | Existing ghost tokens |
| Error | Ghost control may trigger the warning notice; it does not turn red itself | `--color-warning-bg`, `--color-warning-text` for paired notice |
| Empty | May appear as the empty-state next action if a primary action is not required | `--color-text-muted` |

**Accessibility** — Use native anchors for navigation and buttons for actions. Focus ring is required. Minimum hit target should reach 44×44px.

### 2.3 Text Input Field

**Purpose** — Use for entering a todo task name. Do not use for long-form text.

**Anatomy** — `[visually-hidden label] [input] [validation error?]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Standard todo input | `--color-bg`, `--color-primary-border`, `--color-text`, `--radius-lg` | Task entry |
| Error message | `--color-danger` | Blank submission validation |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Default | content height from `15px 16px` padding | `15px 16px` | `--text-base` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Soft background, pale blue border, dark text | `--color-bg`, `--color-primary-border`, `--color-text` |
| Hover | No separate hover style is defined | Standard input tokens |
| Focus (keyboard) | Visible 3px blue translucent focus ring with 3px offset | `--color-focus` |
| Active / pressed | Text caret and native editing behavior; no separate visual style | Standard input tokens |
| Disabled | No disabled style is defined; disabled fields should not be used in the approved flow | Standard input tokens only |
| Loading | Input remains available unless the whole todo list is in loading state | `--color-bg`, `--color-text` |
| Error | Error text appears below after blank submit; input remains focused | `--color-danger` |
| Empty | Placeholder gives example task copy: `Add a task, e.g. Send weekly update` | `--color-placeholder` |

**Accessibility** — The input has a label available to assistive technology. Validation returns focus to the field. Maximum length is 80 characters in the approved mockup.

### 2.4 Todo Item

**Purpose** — Display one task with complete/uncomplete and delete actions.

**Anatomy** — `[completion button] [task text] [delete icon button]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Active todo | `--color-surface`, `--color-border`, `--color-text` | Incomplete tasks |
| Completed todo | `--color-success`, `--color-placeholder` | Completed tasks |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Default | content height from `12px` padding and controls | `12px` | `--text-base` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | White row, gray border, dark task text | `--color-surface`, `--color-border`, `--color-text` |
| Hover | No row hover style is defined; child controls may show hover states | Child control tokens |
| Focus (keyboard) | Focus appears on the checkbox or delete button, not the row container | `--color-focus` |
| Active / pressed | Completion toggles between active and completed visual states | `--color-success`, `--color-placeholder` |
| Disabled | No disabled state is defined; rows are interactive in the approved flow | Existing todo tokens only |
| Loading | Todos are hidden while the loading-state component is visible | `--color-text-muted` |
| Error | Todo row remains unchanged when sync warning appears; user input is preserved | `--color-warning-bg`, `--color-warning-text` for paired notice |
| Empty | The list is replaced by the empty-state component when no rows match | Empty-state tokens |

**Accessibility** — Use a list item inside a live `ul`. Completion and delete are native buttons with `aria-label`. Keyboard users must be able to tab to both controls.

### 2.5 Completion Checkbox Button

**Purpose** — Toggle a todo between active and completed.

**Anatomy** — `[checkmark]` with accessible label.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Incomplete | `--color-surface`, `--color-primary-control` | Task is active |
| Complete | `--color-success`, `--color-primary-text` | Task is done |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Icon | `28px` visual control | none | icon glyph |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | White circle with 2px pale blue border | `--color-surface`, `--color-primary-control`, `--border-width-strong` |
| Hover | No separate hover style is defined | Default tokens |
| Focus (keyboard) | Visible 3px blue translucent focus ring with 3px offset | `--color-focus` |
| Active / pressed | Toggles to green filled circle for completed state | `--color-success`, `--color-primary-text` |
| Disabled | No disabled style is defined | Existing checkbox tokens only |
| Loading | Hidden with todo list during initial loading | Loading-state tokens |
| Error | No error style; sync errors are shown in the notice | Notice tokens |
| Empty | Not displayed when list is empty | Empty-state tokens |

**Accessibility** — Native button with `aria-label` changing between `Mark complete` and `Mark incomplete`. Visual target is 28×28px, below the 44×44px recommended touch target.
### 2.6 Delete Icon Button

**Purpose** — Delete a single todo item. Use only for destructive item-level removal.

**Anatomy** — `[× icon]` with accessible label.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Destructive icon | `--color-danger-bg`, `--color-danger`, `--radius-sm` | Delete todo action |
| Destructive hover | `--color-danger-hover-bg` | Hover state |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Icon | `34px` square | none | icon glyph |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Pale red square with red glyph | `--color-danger-bg`, `--color-danger` |
| Hover | Darker pale red background | `--color-danger-hover-bg` |
| Focus (keyboard) | Visible 3px blue translucent focus ring with 3px offset | `--color-focus` |
| Active / pressed | Deletes the row; no separate pressed visual defined | Destructive icon tokens |
| Disabled | No disabled style is defined | Existing destructive tokens only |
| Loading | Hidden with todo list during initial loading | Loading-state tokens |
| Error | Delete button itself does not show error; sync warning appears separately | Notice tokens |
| Empty | Not displayed when list is empty | Empty-state tokens |

**Accessibility** — Native button with `aria-label="Delete task"`. Visual target is 34×34px, below the 44×44px recommended touch target.

### 2.7 Filter Tabs

**Purpose** — Switch the visible todo list between all, active, and completed tasks.

**Anatomy** — `[tab label]` within `[tablist]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Inactive tab | `--color-surface`, `--color-primary-tint`, `--color-text-muted` | Filter not selected |
| Active tab | `--color-primary-soft`, `--color-primary` | Current filter |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Pill | content height from `9px 12px` padding | `9px 12px` | `--text-base` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | White pill, pale blue border, muted text | `--color-surface`, `--color-primary-tint`, `--color-text-muted` |
| Hover | No separate hover style is defined | Inactive tab tokens |
| Focus (keyboard) | Visible 3px blue translucent focus ring | `--color-focus` |
| Active / pressed | Pale primary background, primary text, heavier font weight | `--color-primary-soft`, `--color-primary` |
| Disabled | No disabled style is defined | Existing tab tokens only |
| Loading | Tabs remain visible while loading state shows below | Loading-state tokens |
| Error | Tabs remain unchanged when sync notice appears | Notice tokens |
| Empty | Active filter may reveal empty-state component when no tasks match | Empty-state tokens |

**Accessibility** — Container uses `role="tablist"` and buttons use `role="tab"`. The approved mockup does not define `aria-selected` management.

### 2.8 Loading State

**Purpose** — Show while existing database-backed todos are fetched.

**Anatomy** — `[three animated dots] [loading message]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Full loading | `--color-text-muted`, `--color-primary-control` | Todo list loading area |
| Mini loading | `--color-primary-control` | State explanation card |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Full | content height from `28px` padding | `28px` | `--text-base` |
| Dot | `10px` | `4px` margin | none |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Three pale blue pulsing dots and muted message | `--color-primary-control`, `--color-text-muted` |
| Hover | No hover state; component is not interactive | Loading tokens |
| Focus (keyboard) | No focus; component is not interactive | Loading tokens |
| Active / pressed | Not interactive | Loading tokens |
| Disabled | Not interactive | Loading tokens |
| Loading | Dots animate with `.8s` alternate pulse and staggered delays | `--duration-loading`, `--easing` |
| Error | Replaced or accompanied by sync notice if fetch fails | Notice tokens |
| Empty | Replaced by empty-state component if there are no tasks after loading | Empty-state tokens |

**Accessibility** — Loading area uses `aria-live="polite"`. Motion should be removed for users who prefer reduced motion, but the approved CSS does not include that media query.

### 2.9 Empty State

**Purpose** — Explain when there are no todos and prompt the user to add one.

**Anatomy** — `[icon] [heading] [supporting text]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Todo empty | `--color-bg`, `--color-primary-border`, `--color-primary-soft`, `--color-primary`, `--color-text`, `--color-text-muted` | Empty todo list |
| Mini empty | `--color-primary`, `--color-primary-control` | State explanation card |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Full | content height from `30px 12px` padding | `30px 12px` | h3 plus body |
| Icon | `54px` square | none | `28px` icon |
| Mini | `86px` | none | `34px` icon |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Dashed pale blue border, soft background, cloud icon tile, heading and muted copy | Empty-state tokens |
| Hover | No hover state; component is not interactive | Empty-state tokens |
| Focus (keyboard) | No focus unless a future action is added | Empty-state tokens |
| Active / pressed | Not interactive | Empty-state tokens |
| Disabled | Not interactive | Empty-state tokens |
| Loading | Hidden while loading-state component is visible | Loading-state tokens |
| Error | Replaced or accompanied by notice if data cannot load | Notice tokens |
| Empty | Message: `No tasks yet` and `Add one small task to get momentum for the day.` | Empty-state tokens |

**Accessibility** — Empty state must not be blank. Keep clear heading and next-action copy. Decorative icon should not be the only information.

### 2.10 Notice / Error Message

**Purpose** — Show sync problems and validation errors without losing user input.

**Anatomy** — `[message]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Inline validation | `--color-danger` | Blank task submission |
| Sync warning notice | `--color-warning-bg`, `--color-warning-text`, `--color-warning-border`, `--radius-md` | Database or sync problem |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Inline | content height | margin `8px 2px 0` | `--text-sm` |
| Notice | content height from `12px 14px` padding | `12px 14px` | `--text-base` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Hidden until needed | Existing state tokens |
| Hover | No hover state; component is not interactive | Existing state tokens |
| Focus (keyboard) | Focus returns to related input for validation; notice itself is not focusable | `--color-focus` on related input |
| Active / pressed | Not interactive | Existing state tokens |
| Disabled | Not interactive | Existing state tokens |
| Loading | Hidden during normal loading unless an error occurs | Loading tokens |
| Error | Inline validation uses red text; sync warning uses amber panel | `--color-danger`, warning tokens |
| Empty | Validation copy explains missing task: `Enter a task before adding it.` | `--color-danger` |

**Accessibility** — Error/notice text must be plain language. Validation should focus the input. Todo list and loading regions use `aria-live="polite"`; future production errors should also be announced.

### 2.11 Card / Panel

**Purpose** — Group related product content, todo UI, states, and details.

**Anatomy** — `[optional header] [content] [optional footer/actions]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| App card | `--color-card-bg`, `--color-border`, `--radius-app`, `--shadow-app` | Main interactive todo preview |
| State card | `--color-surface`, `--color-border`, `--radius-card`, `--shadow-xl` | Loading, empty, and error explanations |
| Stat card | `--color-stat-bg`, `--color-border-soft`, `--radius-xl` | Hero highlight numbers |
| Accordion panel | `--color-bg`, `--color-border`, `--radius-lg` | Product detail answer body |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| App | content-based | `22px` | mixed |
| State | content-based | `20px` | mixed |
| Stat | content-based | `14px 16px` | `--text-sm`, `--text-stat` |
| Panel | content-based | `0 16px` | `--text-base` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Surface, border, radius, and component-specific shadow | Card/panel tokens |
| Hover | No card hover style is defined | Card/panel tokens |
| Focus (keyboard) | Cards are not focusable unless they contain controls | Child control focus tokens |
| Active / pressed | Not interactive as a container | Card/panel tokens |
| Disabled | Not interactive as a container | Card/panel tokens |
| Loading | App card contains loading-state component | Loading tokens |
| Error | App card contains notice/error component | Notice tokens |
| Empty | App card contains empty-state component | Empty-state tokens |

**Accessibility** — Use semantic section/article containers with labelled headings where shown in the approved design.
### 2.12 Accordion

**Purpose** — Reveal short supporting product details without expanding the page by default.

**Anatomy** — `[button label] [plus/minus indicator] [panel]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Closed accordion | `--color-surface`, `--color-primary-tint`, `--radius-lg` | Collapsed detail row |
| Open panel | `--color-bg`, `--color-border`, `--radius-lg` | Expanded detail copy |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Button | content height from `16px` padding | `16px` | `--text-base`, weight `800` |
| Panel | content height | `0 16px` | `--text-base` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | White button, pale blue border, plus indicator | `--color-surface`, `--color-primary-tint` |
| Hover | No separate hover style is defined | Accordion tokens |
| Focus (keyboard) | Visible 3px blue translucent focus ring | `--color-focus` |
| Active / pressed | Toggles panel open and indicator from `+` to `−` | Panel tokens |
| Disabled | No disabled style is defined | Existing accordion tokens only |
| Loading | Not used for loading content | Existing accordion tokens |
| Error | Not used for error content | Existing accordion tokens |
| Empty | Closed state uses no blank body; open panels always include explanatory copy | Panel tokens |

**Accessibility** — Button uses `aria-expanded`. Escape closes open panels in the approved interaction model.

### 2.13 Navigation Header

**Purpose** — Provide sticky access to the app, states, and details sections.

**Anatomy** — `[brand link with logo] [mobile menu button] [navigation links] [primary link]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Desktop sticky nav | `--color-nav-bg`, `--color-nav-border` | Screens wider than `820px` |
| Mobile dropdown nav | `--color-surface-raised`, `--color-primary-tint`, `--radius-xl`, `--shadow-app` | Screens at or below `820px` |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Header inner | content height | `16px 22px` | `--text-base` |
| Mobile dropdown | content height | `10px` | `--text-base` |
| Logo | `36px` square | none | glyph |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Sticky blurred translucent bar with bottom border | `--color-nav-bg`, `--color-nav-border` |
| Hover | Child links use ghost hover state | Ghost tokens |
| Focus (keyboard) | Child controls show focus ring; skip link can move into view above content | `--color-focus` |
| Active / pressed | Mobile menu toggles open/closed and updates `aria-expanded` | Mobile dropdown tokens |
| Disabled | Navigation items are not disabled in the approved design | Existing nav tokens |
| Loading | Navigation remains available while app content loads | Loading-state tokens in app only |
| Error | Navigation remains unchanged when app shows sync notice | Notice tokens in app only |
| Empty | Navigation remains unchanged when todo list is empty | Empty-state tokens in app only |

**Accessibility** — Nav has `aria-label="Main navigation"`. Mobile toggle has `aria-expanded` and `aria-controls`. Escape closes the menu.

### 2.14 Brand Logo

**Purpose** — Identify the product in the header.

**Anatomy** — `[gradient square with checkmark] [product name]`.

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Header brand | `--color-primary`, `--color-primary-gradient-end`, `--color-primary-text`, `--radius-sm`, `--shadow-md` | Main navigation brand |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Logo mark | `36px` square | none | glyph |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Rounded gradient square with white checkmark | Brand tokens |
| Hover | No separate hover style is defined on the brand logo | Brand tokens |
| Focus (keyboard) | Brand link inherits anchor focus behavior if browser default applies; no explicit brand focus selector is defined | Existing brand tokens |
| Active / pressed | Navigates to page top | Brand tokens |
| Disabled | Not disabled | Brand tokens |
| Loading | Unchanged | Brand tokens |
| Error | Unchanged | Brand tokens |
| Empty | Unchanged | Brand tokens |

**Accessibility** — Brand is an anchor to `#top`; the visible product name provides the accessible name.

## 3. Content and formatting

- Voice and tone: calm, plain, focused, and practical; avoid hype and keep copy task-oriented.
- Date, time, number, and currency formats: no dates, times, or currency appear in the approved design; numbers are plain English UI counts such as `0 active tasks`, `1 active task`, and `3 Core actions`.
- Locale: English (`html lang="en"`).
- Capitalization rule for buttons, headings, and labels: sentence case for buttons and headings except the product name `Todo List App v2`; short tab labels use title-style single words where natural (`All`, `Active`, `Completed`).
- Empty-state wording pattern: name what is missing, then give the next small action. Example: `No tasks yet` followed by `Add one small task to get momentum for the day.`
- Error-message wording pattern: plain language, preserve user input, and avoid blaming the user. Examples: `Enter a task before adding it.` and `We could not sync just now. Your typed task is still safe.`

## 4. Known deviations

Places where the approved design does not follow its own rules or the anti-patterns in `references/ai-defaults.md`. Record, do not silently fix.

| Where | Deviation | Why it stands | Follow-up |
|---|---|---|---|
| Body background and logo | Uses gradients: radial page background and linear logo gradient. `references/ai-defaults.md` warns against decorative gradients. | Stakeholder approved a polished blue-and-white visual direction and the approved mockup includes these gradients. | Keep gradients only in these approved locations unless the stakeholder requests flatter styling. |
| Radius scale | Many distinct radii are used: `10px`, `12px`, `14px`, `16px`, `18px`, `20px`, `22px`, `24px`, `999px`, and `9999px`. The default guidance recommends 3-4 radius steps. | The approved mockup uses a broad but consistent soft radius language. | Consolidate radii in implementation only if approved as a design change. |
| Spacing scale | Spacing includes many component-specific values such as `13px`, `15px`, `22px`, `28px`, `34px`, `46px`, `62px`, `70px`, and `80px` rather than a tight 4/8 scale. | These values are present in the approved CSS and define the current layout rhythm. | Treat unusual values as component-specific until the stakeholder approves simplification. |
| Motion accessibility | No `prefers-reduced-motion: reduce` media query exists, while the defaults require reduced-motion support. | The approved mockup includes animations but not the reduced-motion override. | Implementation should add reduced-motion support without changing the visible approved default motion. |
| UI target size | Completion checkbox is `28px` and delete icon button is `34px`, below the recommended 44×44px touch target. | The approved mockup prioritizes compact todo rows. | Increase touch targets only if approved or if implementation can add invisible hit area without changing visuals. |
| Contrast: completed text | Completed todo text uses `#94A3B8` on `#FFFFFF`, about `2.5:1`, below AA for body text. | The approved design uses low contrast to de-emphasize completed tasks. | Adjust completed text contrast only with stakeholder approval. |
| Contrast: delete/error text on pale red | `#EF4444` on `#FEF2F2` is about `3.5:1`, below AA for normal body text. | Used for icon/error emphasis in the approved mockup. | Avoid using this pair for long body copy unless colors are approved for change. |
| Contrast: white on danger | `#FFFFFF` on `#EF4444` is about `3.8:1`, below AA for normal text. | The danger button is present in CSS but not used prominently in the approved visible flow. | Use danger buttons only for large/bold labels or revisit danger color before production. |
| Border contrast | Pale blue checkbox/input/empty borders such as `#93C5FD` and `#BFDBFE` on white fall below 3:1 UI contrast. | The approved mockup uses subtle blue borders for a soft visual style. | Strengthen borders only with stakeholder approval. |
| Tab ARIA | Tabs use `role="tab"` but do not define `aria-selected`. | The approved prototype focuses on visual interaction. | Add `aria-selected` in implementation while preserving visual design. |
| Iconography | Uses glyphs/emoji-like symbols (`✓`, `☰`, `☁`, `+`, `!`, `×`) rather than a formal icon set. `references/ai-defaults.md` warns against emoji as iconography. | Symbols are embedded in the approved HTML and are simple enough for this app. | Replace with a consistent icon set only if the stakeholder requests it. |

## 5. Change log

| Date | Change | Design PR |
|---|---|---|
| 2026-08-11 | Initial design system extracted from approved `index.html`. | This PR |
