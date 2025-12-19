// 侧边栏折叠逻辑
const toggleBtn = document.getElementById('toggleBtn');
const sidebar = document.getElementById('sidebar');

toggleBtn.addEventListener('click', () => {
    sidebar.classList.toggle('aside-collapsed');
    // 如果你希望侧边栏消失后不占位，可以切换 w-64 和 w-0
    if (sidebar.classList.contains('aside-collapsed')) {
        sidebar.style.width = '0';
    } else {
        sidebar.style.width = '16rem'; // 相当于 w-64
    }
});

// 头像下拉菜单切换
const userMenuBtn = document.getElementById('userMenuBtn');
const userDropdown = document.getElementById('userDropdown');

userMenuBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    userDropdown.classList.toggle('hidden');
});

// 点击页面其他地方关闭下拉框
document.addEventListener('click', () => {
    userDropdown.classList.add('hidden');
});