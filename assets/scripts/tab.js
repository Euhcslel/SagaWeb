// Функция для отображения плавильного item у tab
function showItem(e, item) {
  document.querySelectorAll(".tabItem").forEach((content) => {
    content.style.display = "none";
  });

  document.querySelectorAll(".tablinks").forEach((tabLink) => {
    tabLink.classList.remove("active");
  });

  const el = document.getElementById(item);

  if (item === "mainConfig") {
    el.style.display = "flex";
    return;
  }

  el.style.display = "block";
  e.currentTarget.classList.add("active");
}
