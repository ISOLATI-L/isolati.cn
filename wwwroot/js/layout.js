function search(){
    console.log(search_txt.value);
    // alert("未完成");
    return false;
};

const search_txt = document.querySelector(".search_txt")
const search_btn = document.querySelector(".search_btn")

search_btn.onclick = search;
