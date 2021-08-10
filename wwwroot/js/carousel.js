setInterval(function() {
    const show = document.querySelector(".maskText[data-show]");
    const next = show.nextElementSibling ||
    document.querySelector(".maskText:first-child");
    const up = document.querySelector(".maskText[data-up]");

    if (up){
        up.removeAttribute("data-up");
    }
    show.removeAttribute("data-show");
    show.setAttribute("data-up", "");
    next.setAttribute("data-show", "");
}, 2000);
