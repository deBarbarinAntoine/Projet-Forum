const has_list = document.querySelectorAll('.has-list');

function refresh_height() {
    const category_list = document.querySelectorAll('.category-list');
    
    category_list.forEach((list) => {
        list.style.maxHeight = `${list.children.length * 42}px`; 

        // console.log(list.offsetHeight);
        // console.log(list.children.length * 42);
    });
}

has_list.forEach((item) => {
    const item_arrow = item.lastElementChild;
    item.addEventListener('click', () => {
        const list = document.querySelector(`#${item.dataset.list}`);
        list.classList.toggle('closed');
        item_arrow.classList.toggle('rotated');
    });
});

refresh_height();