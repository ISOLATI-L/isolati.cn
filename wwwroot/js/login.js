function loginAdmin() {
    const password = login_txt.value;
    const MD5password = hash(password);
    console.log(MD5password);

    post("/login", MD5password).then(
        function (res) {
            console.log(res);
        },
        function (res) {
            console.log(res);
        });
    return false;
}

function login_txt_keypress(e) {
    if (e.keyCode === 13) {
        return loginAdmin();
    }
    return false;
}

const login_txt = document.querySelector(".login_txt");
const login_btn = document.querySelector(".login_btn");

login_txt.addEventListener("keypress", login_txt_keypress);
login_btn.onclick = loginAdmin;
