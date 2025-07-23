document.addEventListener('DOMContentLoaded', function() {
    let totalPages = 0;
    let data = null;

    // 获取加载相关元素
    const loadingOverlay = document.getElementById('loadingOverlay');
    const loadingLogo = document.querySelector('.loading-logo');
    const loadingText = document.querySelector('.loading-text');
    const loadingProgress = document.querySelector('.loading-progress');
    const body = document.body;

    // 初始状态
    loadingOverlay.style.display = 'flex';
    body.style.overflow = 'hidden';

    fetch('http://localhost:8081/gacha-history')
        .then(response => response.json())
        .then(response => {
            data = response;
            totalPages = Object.keys(data).length;

            // 修改加载文字为"加载完毕"
            const loadingText = document.querySelector('.loading-text');
            loadingText.textContent = '加载完毕';
            loadingText.style.color = '#00796b'; // 保持原有颜色或更改为完成颜色

            // 创建内容
            createCharts(data);
            displayRareChars(data);

            // 优雅的退出动画序列
            startExitAnimation();
        })
        .catch(error => {
            console.error('Error fetching data:', error);
            handleLoadError();
        });

    function startExitAnimation() {
        // 1. 首先完成进度条动画
        loadingProgress.indeterminate = false;
        loadingProgress.value = 100;

        // 2. 同时开始其他元素的退出动画
        setTimeout(() => {
            // 文字淡出
            loadingText.style.transition = 'all 0.4s ease-out';
            loadingText.style.opacity = '0';
            loadingText.style.transform = 'translateY(-10px)';

            // 进度条淡出
            loadingProgress.style.transition = 'opacity 0.3s ease-out';
            loadingProgress.style.opacity = '0';

            // logo缩放和淡出
            loadingLogo.style.transition = 'all 0.5s cubic-bezier(0.68, -0.55, 0.27, 1.55)';
            loadingLogo.style.transform = 'scale(0.8) translateY(-20px)';
            loadingLogo.style.opacity = '0';

            // 3. 整体遮罩上滑消失
            setTimeout(() => {
                loadingOverlay.style.transition = 'all 0.6s cubic-bezier(0.64, 0, 0.35, 1)';
                loadingOverlay.style.transform = 'translateY(-100%)';
                loadingOverlay.style.opacity = '0';

                // 4. 恢复页面状态
                setTimeout(() => {
                    body.style.overflow = 'auto';
                    loadingOverlay.style.display = 'none';

                    // 5. 内容轻微弹入效果
                    const contentElements = document.querySelectorAll('h1, .chart-container, .rare-chars');
                    contentElements.forEach((el, index) => {
                        el.style.opacity = '0';
                        el.style.transform = 'translateY(20px)';
                        el.style.transition = `all 0.5s ease-out ${index * 0.1}s`;
                        setTimeout(() => {
                            el.style.opacity = '1';
                            el.style.transform = 'translateY(0)';
                        }, 50);
                    });
                }, 600);
            }, 300);
        }, 200);
    }

    function handleLoadError() {
        // 修改文字提示
        loadingText.textContent = '加载失败，请刷新重试';
        loadingText.style.color = '#ff4444';

        // 进度条变红色
        loadingProgress.style.setProperty('--mdui-color-primary', '255, 68, 68');

        // 错误状态动画
        setTimeout(() => {
            // 震动效果提示错误
            loadingOverlay.style.animation = 'shake 0.5s cubic-bezier(.36,.07,.19,.97) both';

            // 3秒后淡出
            setTimeout(() => {
                loadingOverlay.style.transition = 'all 0.8s ease';
                loadingOverlay.style.opacity = '0';

                setTimeout(() => {
                    loadingOverlay.style.display = 'none';
                    body.style.overflow = 'auto';
                }, 800);
            }, 3000);
        }, 500);
    }

    function createCharts(data) {
        const chartContainer = document.getElementById('chartContainer');
        Object.keys(data)        // 所有池名
            .slice(-4)           // 取最后 4 个
            .reverse()           // 倒序
            .forEach(poolName => {
            const items = data[poolName];
            const rarityCounts = {2: 0, 3: 0, 4: 0, 5: 0}; // 对应三星、四星、五星、六星
            const rareFiveChars = [];

            items.forEach(item => {
                rarityCounts[item.rarity] += 1;
                if (item.rarity === 5) {
                    rareFiveChars.push(item.charName);
                    if (rareFiveChars.length > 5) {
                        rareFiveChars.shift();
                    }
                }
            });

            const ctx = document.createElement('canvas');
            const canvasId = `chart_${poolName}`;
            ctx.id = canvasId;
            chartContainer.appendChild(ctx);

            new Chart(ctx, {
                type: 'pie',
                data: {
                    labels: ['3★', '4★', '5★', '6★'], // 对应三星、四星、五星、六星
                    datasets: [{
                        label: poolName,
                        data: Object.values(rarityCounts),
                        backgroundColor: ['#ff6384', '#36a2eb', '#cc65fe', '#ffce56'],
                        hoverOffset: 4
                    }]
                },
                options: {
                    responsive: false, // 设置为false以保持固定大小
                    plugins: {
                        legend: {
                            position: 'top'
                        },
                        title: {
                            display: true,
                            text: `${poolName} 寻访记录`,
                            font: {
                                size: 20
                            }
                        },
                    }
                }
            });
        });
    }

    function displayRareChars(data) {
        const rareCharsContainer = document.getElementById('rareCharsContainer');
        rareCharsContainer.style.display = 'flex';
        rareCharsContainer.style.gap = '12px';
        rareCharsContainer.style.flexWrap = 'wrap';

        Object.keys(data)
            .slice(-4)
            .reverse()
            .forEach(poolName => {
                const items = data[poolName];
                const rareFiveChars = [];
                const count = items.length;

                items.forEach(item => {
                    if (item.rarity === 5) {
                        rareFiveChars.push(item.charName);
                        if (rareFiveChars.length > 5) rareFiveChars.shift();
                    }
                });

                const card = document.createElement('mdui-card');
                card.setAttribute('variant', 'elevated');
                card.style.cssText = 'width:300px; margin-left:15px; margin-right:15px; min-height:224px; background:transparent;';
                card.innerHTML = `
                <div style="padding:16px; display:flex; flex-direction:column; gap:12px;">
                    <div style="font-size:18px; font-weight:600; color:#00796b; width: 240px">
                        《${poolName}》
                        <span style="font-size:14px; font-weight:600; color:#00796b; margin-left:4px;">
                            总寻访数：${count}
                        </span>
                    </div>
                    <!-- 干员列表 -->
                    <div style="display:flex; gap:6px; margin-left: 10px">
                        ${rareFiveChars.map(name =>
                            `<mdui-chip style="
                                        background-color:transparent;
                                        color:#004d40;
                                        font-size:14px;
                                        height:28px;
                                        --mdui-chip-elevation:0 1px 3px rgba(0,0,0,.12);
                                    ">${name}</mdui-chip>`
                    ).join('')}
                    </div>
                </div>
            `;

                rareCharsContainer.appendChild(card);
            });
    }
});