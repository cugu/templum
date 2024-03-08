document.addEventListener('DOMContentLoaded', () => {
    document.documentElement.classList.remove('no-js');

    const menu_toggles = document.querySelectorAll('.menu-toggle');
    menu_toggles.forEach(menu_toggle => {
        menu_toggle.addEventListener('click', () => {
            toggleSidebar();
        });
    });

    const overlay = document.getElementById('overlay');
    if (overlay !== null) {
        overlay.addEventListener('click', () => {
            if (overlay.classList.contains('opacity-70')) {
                toggleSidebar();
            }
        });
    }
});

function toggleSidebar() {
    let menu = document.getElementById('sidebar');
    menu.classList.toggle('open');

    toggleOverlay();

    let body = document.getElementsByTagName('body')[0];
    body.classList.toggle('overflow-hidden');
}

function toggleOverlay() {
    let sidebarOverlay = document.getElementById('overlay');
    sidebarOverlay.classList.toggle('opacity-0');
    sidebarOverlay.classList.toggle('opacity-70');
    setTimeout(() => {
        sidebarOverlay.classList.toggle('hidden');
    }, 200);
}
