document.querySelector("form").addEventListener("submit", async function (e) {
  e.preventDefault();
  const response = await fetch("/sign_in", {
    method: "POST",
    body: new URLSearchParams(new FormData(e.target)),
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
  });
  if (response.ok) {
    window.location.href = response.url;
    return;
  }
  const text = await response.text();
  Swal.fire({ icon: "error", title: "Ошибка", text });
});
