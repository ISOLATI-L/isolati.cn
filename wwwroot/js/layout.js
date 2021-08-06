window.onload = function(){
    let transition_noneItems=
        document
        .querySelectorAll(".transition_init");
    for(let i = 0; i < transition_noneItems.length; i++){
        transition_noneItems[i].classList.remove("transition_init");
    }
};
