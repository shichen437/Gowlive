<template>
    <div class="w-full h-full">
        <div ref="playerContainer" class="w-full h-full"></div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue';
import Artplayer from 'artplayer';
import mpegts from 'mpegts.js';

type MediaFormat = 'flv' | 'mp4' | 'mp3' | 'mkv' | 'ts' | 'm3u8';

const props = defineProps<{
    url: string;
    format: MediaFormat;
    isLive: boolean;
    headers?: Record<string, string>;
}>();

const playerContainer = ref<HTMLElement | null>(null);
let player: Artplayer | null = null;
let flvPlayer: mpegts.Player | null = null;
let tsPlayer: mpegts.Player | null = null;

function initializePlayer() {
    const container = playerContainer.value;

    if (!container) {
        console.warn('VideoPlayer container is not mounted yet.');
        return;
    }

    if (!props.url) {
        console.warn('No media url provided to VideoPlayer.');
        destroyPlayer();
        return;
    }

    destroyPlayer();

    const options: any = getBasicOptions(container);

    if (props.format === 'ts') {
        options.type = 'ts';
        options.customType = createTsCustomType();
    }
    if (props.format === 'flv') {
        options.type = 'flv';
        options.customType = createFlvCustomType();
    } else if (props.format === 'm3u8') {
        options.type = 'm3u8';
        if (props.headers) {
            options.hlsOption = {
                xhrSetup: function (xhr: XMLHttpRequest) {
                    for (const key in props.headers) {
                        if (Object.prototype.hasOwnProperty.call(props.headers, key)) {
                            xhr.setRequestHeader(key, props.headers[key]);
                        }
                    }
                }
            };
        }
    } else {
        options.id = options.url
        options.autoPlayback = true
    }
    if (props.format === 'mp3') {
        options.fullscreen = false;
        options.fullscreenWeb = false;
        options.pip = false;
        options.setting = false;
        options.screenshot = false;
    }

    player = new Artplayer(options);
    player.contextmenu.show = false;
    player.on('ready', () => {
        console.log('Player is ready!')
    })
    player.on('error', (error: unknown) => {
        console.error('Artplayer error:', error);
    });
}

function createTsCustomType() {
    return {
        ts: (video: HTMLMediaElement, url: string, art: Artplayer) => {
            if (!mpegts.isSupported() || !mpegts.getFeatureList().mseLivePlayback) {
                art.notice.show = ('当前环境不支持 TS 播放');
                return;
            }
            const mediaDataSource: mpegts.MediaDataSource = { type: 'mpegts', url, isLive: props.isLive };
            if (props.headers) {
                mediaDataSource.headers = props.headers;
            }
            const config: mpegts.Config = {
                isLive: props.isLive,
                enableStashBuffer: !props.isLive,
                stashInitialSize: 512,
                autoCleanupSourceBuffer: !props.isLive,
                accurateSeek: !props.isLive,
                seekType: 'range',
            };
            tsPlayer = mpegts.createPlayer(mediaDataSource, config);
            tsPlayer.attachMediaElement(video);
            tsPlayer.load();
            if (props.isLive) video.play().catch(() => { });
            art.on('destroy', () => destroyTsPlayer());
        },
    };
}

function createFlvCustomType() {
    return {
        flv: (video: HTMLMediaElement, url: string, art: Artplayer) => {
            destroyFlvPlayer();

            if (!mpegts.isSupported()) {
                art.notice.show = ('当前环境不支持 FLV 播放');
                return;
            }

            const mediaDataSource: mpegts.MediaDataSource = {
                type: 'flv',
                url,
                isLive: props.isLive,
            };
            if (props.headers) {
                mediaDataSource.headers = props.headers;
            }

            const config: mpegts.Config = {
                isLive: props.isLive,
                enableStashBuffer: !props.isLive,
                stashInitialSize: 512,
                enableWorker: false,
                autoCleanupSourceBuffer: !props.isLive,
                accurateSeek: !props.isLive,
                seekType: 'range',
                lazyLoad: !props.isLive,
                lazyLoadMaxDuration: 60,
                lazyLoadRecoverDuration: 30,
                statisticsInfoReportInterval: 1000
            };

            flvPlayer = mpegts.createPlayer(mediaDataSource, config);

            flvPlayer.on(mpegts.Events.ERROR, (type, detail) => {
                console.error('mpegts error:', type, detail);
                art.notice.show = (`FLV 播放错误: ${type}`);
            });

            flvPlayer.attachMediaElement(video);
            flvPlayer.load();
            const onSeeking = () => {
                try {
                    if (!video.buffered || video.buffered.length === 0) return;

                    const end = video.buffered.end(video.buffered.length - 1);

                    if (video.currentTime >= end) {
                        const safetyBack = 1;
                        const target = Math.max(0, end - safetyBack);
                        video.currentTime = target;
                    }
                } catch (e) {
                    console.warn('seeking adjust error:', e);
                }
            };

            video.addEventListener('seeking', onSeeking);

            art.flv = flvPlayer;

            art.on('destroy', () => {
                video.removeEventListener('seeking', onSeeking);
                destroyFlvPlayer();
            });
        },
    };
}

function destroyFlvPlayer() {
    if (flvPlayer) {
        try {
            flvPlayer.unload();
            flvPlayer.detachMediaElement();
            flvPlayer.destroy();
        } catch (e) {
            console.warn('mpegts destroy error:', e);
        } finally {
            flvPlayer = null;
        }
    }
}

function destroyTsPlayer() {
    if (tsPlayer) {
        try {
            tsPlayer.unload();
            tsPlayer.detachMediaElement();
            tsPlayer.destroy();
        } catch (e) {
            console.warn('mpegts destroy error:', e);
        } finally {
            tsPlayer = null;
        }
    }
}

function destroyPlayer() {
    destroyFlvPlayer();
    if (player) {
        try {
            player.destroy();
        } catch (e) {
            console.warn('Artplayer destroy error:', e);
        } finally {
            player = null;
        }
    }
}

function getBasicOptions(container: HTMLElement | null) {
    const options: any = {
        airplay: true,
        aspectRatio: true,
        autoOrientation: true,
        autoplay: props.isLive,
        container,
        customType: {},
        flip: true,
        fullscreen: true,
        fullscreenWeb: true,
        lang: 'zh-cn',
        isLive: props.isLive,
        miniProgressBar: true,
        muted: false,
        mutex: true,
        pip: true,
        playbackRate: true,
        poster: '',
        setting: true,
        screenshot: true,
        theme: '#de7897',
        url: props.url,
        volume: 1,
    };
    return options;
}

onMounted(() => {
    initializePlayer();
});

onBeforeUnmount(() => {
    destroyPlayer();
});

watch(
    () => [props.url, props.format, props.isLive, props.headers],
    ([newUrl, newFormat, newIsLive, newHeaders], [oldUrl, oldFormat, oldIsLive, oldHeaders]) => {
        if (
            newUrl !== oldUrl ||
            newFormat !== oldFormat ||
            newIsLive !== oldIsLive ||
            JSON.stringify(newHeaders) !== JSON.stringify(oldHeaders)
        ) {
            initializePlayer();
        }
    }
);

function getCurrentTime() {
    return player ? player.currentTime : 0;
}

function getDuration() {
    return player ? player.duration : 0;
}

function play() {
    player?.play();
}

function pause() {
    player?.pause();
}

defineExpose({
    getCurrentTime,
    getDuration,
    play,
    pause
});
</script>
