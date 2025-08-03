document.addEventListener('DOMContentLoaded', function() {
    let data = null;
    let currentPool = null; // 当前选中的池子

    // 获取加载相关元素
    const loadingOverlay = document.getElementById('loadingOverlay');
    const loadingLogo = document.querySelector('.loading-logo');
    const loadingText = document.querySelector('.loading-text');
    const loadingProgress = document.querySelector('.loading-progress');
    const body = document.body;
    const poolSelector = document.getElementById('poolSelector');

    // 初始状态
    loadingOverlay.style.display = 'flex';
    body.style.overflow = 'hidden';

    fetch('http://localhost:8081/gacha-history')
        .then(response => response.json())
        .then(response => {
            data = response;

            // 修改加载文字为"加载完毕"
            loadingText.textContent = '加载完毕';
            loadingText.style.color = '#00796b';

            // 创建池子选择按钮
            createPoolButtons();

            // 默认显示第一个池子
            if (Object.keys(data).length > 0) {
                currentPool = Object.keys(data)[0];
                updateDisplay();
            }

            // 优雅的退出动画序列
            startExitAnimation();
        })
        .catch(error => {
            console.error('Error fetching data:', error);
            handleLoadError();
        });

    function startExitAnimation() {
        setTimeout(() => {
            loadingText.style.transition = 'all 0.4s ease-out';
            loadingText.style.opacity = '0';
            loadingText.style.transform = 'translateY(-10px)';

            loadingProgress.style.transition = 'opacity 0.3s ease-out';
            loadingProgress.style.opacity = '0';

            loadingLogo.style.transition = 'all 0.5s cubic-bezier(0.68, -0.55, 0.27, 1.55)';
            loadingLogo.style.transform = 'scale(0.8) translateY(-20px)';
            loadingLogo.style.opacity = '0';

            setTimeout(() => {
                loadingOverlay.style.transition = 'all 0.6s cubic-bezier(0.64, 0, 0.35, 1)';
                loadingOverlay.style.transform = 'translateY(-100%)';
                loadingOverlay.style.opacity = '0';

                setTimeout(() => {
                    body.style.overflow = 'auto';
                    loadingOverlay.style.display = 'none';

                    const contentElements = document.querySelectorAll('h1, .pool-selector, .chart-container, .rare-chars');
                    contentElements.forEach((el, index) => {
                        el.style.opacity = '0';
                        el.style.transform = 'translateY(20px)';
                        el.style.transition = `all 0.5s ease-out ${index * 0.1}s`;
                        setTimeout(() => {
                            el.style.opacity = '1';
                            el.style.transform = 'translateY(0)';
                        }, 50);
                    });
                }, 100);
            }, 300);
        }, 200);
    }

    function handleLoadError() {
        loadingText.textContent = '加载失败，请刷新重试';
        loadingText.style.color = '#ff4444';
        loadingProgress.style.setProperty('--mdui-color-primary', '255, 68, 68');

        setTimeout(() => {
            loadingOverlay.style.animation = 'shake 0.5s cubic-bezier(.36,.07,.19,.97) both';

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

    function createPoolButtons() {
        // 清空现有按钮
        poolSelector.innerHTML = '';

        // 创建按钮组
        const buttonGroup = document.createElement('mdui-segmented-button-group');

        // 为每个池子创建按钮
        Object.keys(data).forEach((poolName, index) => {
            const button = document.createElement('mdui-segmented-button');
            button.textContent = poolName;
            button.dataset.poolName = poolName;

            // 第一个按钮默认选中
            if (index === 0) {
                button.setAttribute('selected', '');
            }

            button.addEventListener('click', () => {
                currentPool = poolName;
                updateDisplay();
            });

            buttonGroup.appendChild(button);
        });

        poolSelector.appendChild(buttonGroup);
    }

    function updateDisplay() {
        // 清空现有内容
        const chartContainer = document.getElementById('chartContainer');
        const rareCharsContainer = document.getElementById('rareCharsContainer');
        chartContainer.innerHTML = '';
        rareCharsContainer.innerHTML = '';

        if (!currentPool || !data[currentPool]) return;

        // 创建图表
        createChart(currentPool);

        // 创建稀有干员卡片
        createRareCharsCard(currentPool);
    }

    function createChart(poolName) {
        const items = data[poolName];
        const rarityCounts = {2: 0, 3: 0, 4: 0, 5: 0};
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

        const chartContainer = document.getElementById('chartContainer');
        const ctx = document.createElement('canvas');
        ctx.id = `chart_${poolName}`;
        chartContainer.appendChild(ctx);

        new Chart(ctx, {
            type: 'pie',
            data: {
                labels: ['3★', '4★', '5★', '6★'],
                datasets: [{
                    label: poolName,
                    data: Object.values(rarityCounts),
                    backgroundColor: ['#ff6384', '#36a2eb', '#cc65fe', '#ffce56'],
                    hoverOffset: 4
                }]
            },
            options: {
                responsive: false,
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
    }

    function createRareCharsCard(poolName) {
        const items = data[poolName];
        const rareCharsContainer = document.getElementById('rareCharsContainer');
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
                <div style="display:flex; gap:6px; margin-left: 10px; margin-right: 10px; flex-wrap: wrap; overflow:auto; width: 260px;">
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
    }
});