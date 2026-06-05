(() => {
    const display = document.getElementById('display');
    const status = document.getElementById('status');
    const btnStart = document.getElementById('btn-start');
    const btnStop = document.getElementById('btn-stop');
    const btnReset = document.getElementById('btn-reset');
    const presetBtns = document.querySelectorAll('.preset-btn');

    // State: duration is what we reset to; remainingMs is what's on the clock.
    let durationMs = 300_000; // default S0: 5:00
    let remainingMs = durationMs;
    let running = false;
    let lastTick = 0;
    let rafId = 0;

    function format(ms) {
        const totalSec = Math.max(0, Math.ceil(ms / 1000));
        const m = Math.floor(totalSec / 60);
        const s = totalSec % 60;
        const pad = (n) => String(n).padStart(2, '0');
        return `${pad(m)}:${pad(s)}`;
    }

    function render() {
        display.textContent = format(remainingMs);
        const finished = remainingMs <= 0 && durationMs > 0 && !running;
        display.classList.toggle('finished', finished);
        if (finished) {
            status.textContent = 'Finished';
            document.title = '00:00 - Finished';
        } else if (running) {
            status.textContent = 'Running';
            document.title = `${format(remainingMs)} - Running`;
        } else if (remainingMs === durationMs) {
            status.textContent = 'Ready';
            document.title = 'go-timer';
        } else {
            status.textContent = 'Paused';
            document.title = `${format(remainingMs)} - Paused`;
        }
    }

    function tick(now) {
        if (!running) return;
        const dt = now - lastTick;
        lastTick = now;
        remainingMs -= dt;
        if (remainingMs <= 0) {
            remainingMs = 0;
            running = false;
            render();
            try { beep(); } catch (_) { /* audio is best-effort */ }
            return;
        }
        render();
        rafId = requestAnimationFrame(tick);
    }

    function start() {
        if (running || remainingMs <= 0) return;
        running = true;
        lastTick = performance.now();
        rafId = requestAnimationFrame(tick);
        render();
    }

    function stop() {
        if (!running) return;
        running = false;
        cancelAnimationFrame(rafId);
        render();
    }

    function reset() {
        running = false;
        cancelAnimationFrame(rafId);
        remainingMs = durationMs;
        render();
    }

    function setDuration(seconds) {
        const s = Math.max(0, Math.min(24 * 60 * 60, seconds | 0));
        running = false;
        cancelAnimationFrame(rafId);
        durationMs = s * 1000;
        remainingMs = durationMs;
        render();
    }

    // Small audible cue when finished. Uses WebAudio so we don't need an asset.
    let audioCtx;
    function beep() {
        audioCtx = audioCtx || new (window.AudioContext || window.webkitAudioContext)();
        const o = audioCtx.createOscillator();
        const g = audioCtx.createGain();
        o.type = 'sine';
        o.frequency.value = 880;
        o.connect(g);
        g.connect(audioCtx.destination);
        g.gain.setValueAtTime(0.001, audioCtx.currentTime);
        g.gain.exponentialRampToValueAtTime(0.2, audioCtx.currentTime + 0.02);
        g.gain.exponentialRampToValueAtTime(0.001, audioCtx.currentTime + 0.6);
        o.start();
        o.stop(audioCtx.currentTime + 0.6);
    }

    btnStart.addEventListener('click', start);
    btnStop.addEventListener('click', stop);
    btnReset.addEventListener('click', reset);

    presetBtns.forEach((btn) => {
        btn.addEventListener('click', () => {
            presetBtns.forEach((b) => b.classList.remove('active'));
            btn.classList.add('active');
            setDuration(parseInt(btn.dataset.seconds, 10));
        });
    });

    // Keyboard shortcuts: space = start/stop, r = reset.
    document.addEventListener('keydown', (ev) => {
        if (ev.code === 'Space') { ev.preventDefault(); running ? stop() : start(); }
        else if (ev.key === 'r' || ev.key === 'R') { reset(); }
    });

    render();
})();
