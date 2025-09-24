export function getGreeting(nickname: string) {
    const now = new Date();
    const hour = now.getHours();
    if (hour >= 5 && hour < 11) {
        return '早上好，' + nickname + '！今天也要元气满满哦！';
    }
    if (hour >= 11 && hour < 13) {
        return '中午好，' + nickname + '！吃饱喝足，午休一会儿更舒服！';
    }
    if (hour >= 13 && hour < 19) {
        return '下午好，' + nickname + '！继续加油，保持好心情！';
    }
    if (hour >= 19 && hour < 23) {
        return '晚上好，' + nickname + '！辛苦一天，放松一下吧！';
    }
    return '夜深了，盖好被子，做个好梦，晚安~';
}