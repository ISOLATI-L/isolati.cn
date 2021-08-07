function search(){
    console.log(search.search_txt.value);
    // alert("未完成");
    return false;
};

search.search_txt = document.querySelector(".search_txt")
document.querySelector(".search_btn").onclick = search;
