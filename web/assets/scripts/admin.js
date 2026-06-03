function openUpdateModal(tableName, button) {
  const tableRow = button.closest(".table-row");
  const modal = document.getElementById("edit-modal");

  switch (tableName) {
    case "colors":
      const hex = tableRow.dataset.hex;
      const code = tableRow.dataset.code;
      const id = tableRow.dataset.id;

      const inputId = modal.querySelector(".id");
      inputId.value = id;
      const inputHex = modal.querySelector(".hex");
      inputHex.value = hex;
      const inputCode = modal.querySelector(".code");
      inputCode.value = code;
    case "cycle_amounts":
    case "industrial_drives":
    case "kift_types":
    case "options":
    case "products":
    case "residential_drives":
  }

  modal.showModal();
}

async function updateRow(event) {
  event.preventDefault();

  const form = event.target;
  const url = form.dataset.url;

  const formData = new FormData(form);

  const response = await fetch(url, {
    method: "PUT",
    body: new URLSearchParams(formData),
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
  });

  if (!response.ok) {
    return;
  }

  window.location.reload();
}

async function deleteRow(rowId, tableName) {
  const response = await fetch('/tables/'+tableName+'/'+rowId, {
    method: "DELETE",
  });

  if (!response.ok) {
    return;
  }

  window.location.reload();
}
