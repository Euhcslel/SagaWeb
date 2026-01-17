const protoFile = `
    syntax = "proto3";
    message User {
        int64 id = 1;
        string userName = 2;
        string fullName = 3;
        string email = 4;
        string phoneNumber = 5;
    }
  `;
async function loadUserData() {
  try {
    const root = protobuf.parse(protoFile, { keepCase: true }).root;
    const User = root.lookupType("User");
    const userId = window.location.pathname.split("/").pop();

    const response = await fetch(`/api/users/${userId}`);
    if (!response.ok) throw new Error(`HTTP ${response.status}`);

    const buffer = await response.arrayBuffer();
    const message = User.decode(new Uint8Array(buffer));
    const user = User.toObject(message, { longs: String, keepCase: true });
    document.querySelector(".full-name").textContent = user.fullName;
    document.querySelector(".user-id").textContent = user.id;
    document.querySelector(".username").textContent = user.userName;
    document.querySelector(".email").textContent = user.email;
    document.querySelector(".phone").textContent = user.phoneNumber;
  } catch (error) {
    console.error("Error:", error);
    document.getElementById(
      "user-info"
    ).innerHTML = `<p class="error">Ошибка загрузки: ${error.message}</p>`;
  }
}
document.addEventListener("DOMContentLoaded", loadUserData);
