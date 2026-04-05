window.deleteGate = async (gateId) => {
  const path = window.location.pathname;
  const parts = path.split("/");
  const saleId = parts[2];

  const res = await fetch("/orders/" + saleId + "/" + gateId, {
    method: "DELETE",
  });

  if (!res.ok) throw new Error("Delete request failed");

  window.location.reload();
};

window.deleteOrder = async (saleId) => {
  const res = await fetch("/orders/" + saleId, {
    method: "DELETE",
  });

  if (!res.ok) throw new Error("Delete request failed");

  window.location.href = "/orders";
};

window.addGate = async (gateType) => {
  if (!gateType) return;

  const path = window.location.pathname;

  const formData = new FormData();
  formData.append("gateType", gateType);

  try {
    const response = await fetch(path, {
      method: "POST",
      body: formData,
    });

    if (response.status === 201) {
      const redirect = response.headers.get("Location");
      if (redirect) {
        window.location.href = path + redirect;
      }
    }
  } catch (error) {
    console.error("Сетевая ошибка:", error);
  }
};
