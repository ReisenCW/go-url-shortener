function showError(message) {
    const errorEl = document.getElementById('errorMsg');
    errorEl.textContent = message;
    errorEl.classList.add('visible');
    
    // 3秒后自动隐藏错误提示
    setTimeout(() => {
        errorEl.classList.remove('visible');
    }, 3000);
}

function hideError() {
    const errorEl = document.getElementById('errorMsg');
    errorEl.classList.remove('visible');
}

function onSubmit() {
    const urlInput = document.getElementById('urlInput');
    const userIdInput = document.getElementById('userIdInput');
    const shortenedInput = document.getElementById('shortenedUrl');
    const shortenButton = document.getElementById('shortenBtn');
    const resultCard = document.getElementById('resultCard');

    const longUrl = urlInput.value.trim();
    const userId = userIdInput.value.trim() || 'anonymous';

    // 验证URL格式
    if (!longUrl) {
        showError('请输入需要缩短的URL');
        return;
    }

    // 简单的URL格式验证
    try {
        new URL(longUrl);
    } catch (e) {
        showError('请输入有效的URL（例如：https://example.com）');
        return;
    }

    // 重置状态
    hideError();
    resultCard.classList.remove('visible');
    shortenButton.disabled = true;
    shortenButton.classList.add('loading');
    shortenedInput.value = '';

    console.log('Sending request to shorten URL:', longUrl, 'for user:', userId);

    fetch('http://localhost:9808/create-short-url', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ long_url: longUrl, user_id: userId })
    })
    .then(response => {
        if (!response.ok) throw new Error(`服务器返回错误: ${response.status}`);
        return response.json();
    })
    .then(data => {
        if (data.short_url) {
            shortenedInput.value = data.short_url;
            resultCard.classList.add('visible');
            
            // 短链接点击事件
            shortenedInput.addEventListener('click', () => {
                const href = shortenedInput.value;
                if (href) window.open(href, '_blank', 'noopener');
            });
        } else if (data.error) {
            showError(`错误: ${data.error}`);
        } else {
            showError('服务器返回异常响应');
        }
    })
    .catch(err => {
        showError(`请求失败: ${err.message}`);
    })
    .finally(() => {
        shortenButton.disabled = false;
        shortenButton.classList.remove('loading');
    });
}

// 初始化事件监听
document.addEventListener('DOMContentLoaded', () => {
    // 回车提交
    const urlInput = document.getElementById('urlInput');
    urlInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            onSubmit();
        }
    });

    // 复制按钮功能
    const copyBtn = document.getElementById('copyBtn');
    copyBtn.addEventListener('click', () => {
        const shortVal = document.getElementById('shortenedUrl').value;
        if (!shortVal) return;

        navigator.clipboard.writeText(shortVal)
            .then(() => {
                const originalText = copyBtn.innerHTML;
                copyBtn.innerHTML = '<i class="fas fa-check"></i><span>已复制</span>';
                
                setTimeout(() => {
                    copyBtn.innerHTML = originalText;
                }, 1500);
            })
            .catch(() => {
                showError('复制失败，请手动复制');
            });
    });
});