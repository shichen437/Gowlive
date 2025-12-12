export function getGreetingTime() {
    const now = new Date();
    const hour = now.getHours();
    if (hour >= 5 && hour < 11) {
        return 'morning';
    }
    if (hour >= 11 && hour < 13) {
        return 'noon';
    }
    if (hour >= 13 && hour < 19) {
        return 'afternoon';
    }
    if (hour >= 19 && hour < 23) {
        return 'evening';
    }
    return 'night';
}