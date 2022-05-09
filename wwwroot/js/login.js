function loginAdmin() {
    const password = login_txt.value;
    const MD5password = hash(password);
    console.log(MD5password);

    post("/login", MD5password, "application/x-www-form-urlencoded").then(
        function (res) {
            console.log(res);
            if (res.status === 200) {
                let referer = window.location.search.substring(1).match("ref=(\"|%22)(.*)(\"|%22)(&|$)")
                if (referer.length > 4 && referer[2].length > 0) {
                    referer = referer[2].replaceAll("-", "+")
                    referer = referer.replaceAll("_", "/")
                    window.location.href = window.atob(referer)
                } else {
                    window.location.href = '/admin'
                }
            } else if (res.status === 401) {
                alert("密码错误！");
            } else {
                alert("登陆失败！");
            }
        },
        function (res) {
            console.log(res);
            alert("登陆失败！");
        });
    return false;
}

function login_txt_keypress(e) {
    if (e.keyCode === 13) {
        loginAdmin()
    }
}

const login_txt = document.querySelector(".login_txt");
const login_btn = document.querySelector(".login_btn");

login_txt.addEventListener("keypress", login_txt_keypress);
login_btn.onclick = loginAdmin;
