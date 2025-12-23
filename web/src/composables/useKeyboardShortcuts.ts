import { onMounted, onUnmounted } from 'vue';

type ShortcutCallback = (event: KeyboardEvent) => void;
type Shortcuts = Record<string, ShortcutCallback>;

export function useKeyboardShortcuts(shortcuts: Shortcuts) {
    const handleKeyDown = (event: KeyboardEvent) => {
        const key = event.key.toLowerCase();
        const shortcut = Object.keys(shortcuts).find(s => {
            const parts = s.toLowerCase().split('+');
            const eventKey = parts.pop();
            const modifiers = parts;

            if (key !== eventKey) {
                return false;
            }

            const ctrl = modifiers.includes('ctrl');
            const shift = modifiers.includes('shift');
            const alt = modifiers.includes('alt');
            const meta = modifiers.includes('meta');

            if (event.ctrlKey !== ctrl || event.shiftKey !== shift || event.altKey !== alt || event.metaKey !== meta) {
                return false;
            }

            return true;
        });

        if (shortcut) {
            event.preventDefault();
            shortcuts[shortcut](event);
        }
    };

    onMounted(() => {
        window.addEventListener('keydown', handleKeyDown);
    });

    onUnmounted(() => {
        window.removeEventListener('keydown', handleKeyDown);
    });
}
