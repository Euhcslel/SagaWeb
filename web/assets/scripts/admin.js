function openUpdateModal(tableName, button) {
  const tableRow = button.closest(".table-row");
  const modal = document.getElementById("edit-modal");

  switch (tableName) {
    case "colors":
      modal.querySelector(".id").value = tableRow.dataset.id;
      modal.querySelector(".code").value = tableRow.dataset.code;
      modal.querySelector(".hex").value = tableRow.dataset.hex;
      break;
    case "cycle_amounts":
      modal.querySelector(".id").value = tableRow.dataset.id;
      modal.querySelector(".amount").value = tableRow.dataset.amount;
      modal.querySelector(".wholesale-markup").value =
        tableRow.dataset.wholesaleMarkup;
      modal.querySelector(".retail-markup").value =
        tableRow.dataset.retailMarkup;
      break;
    case "industrial_drives":
      modal.querySelector(".id").value = tableRow.dataset.id;
      modal.querySelector(".name").value = tableRow.dataset.name;
      modal.querySelector(".wholesale-price").value =
        tableRow.dataset.wholesalePrice;
      modal.querySelector(".retail-price").value = tableRow.dataset.retailPrice;
      modal.querySelector(".specifications").value =
        tableRow.dataset.specifications || "";
      modal.querySelector(".image-path").value =
        tableRow.dataset.imagePath || "no_image";
      break;
    case "lift_types":
      modal.querySelector(".id").value = tableRow.dataset.id;
      modal.querySelector(".name").value = tableRow.dataset.name;
      modal.querySelector(".min-headroom").value = tableRow.dataset.minHeadroom;
      modal.querySelector(".max-headroom").value = tableRow.dataset.maxHeadroom;
      modal.querySelector(".wholesale-markup").value =
        tableRow.dataset.wholesaleMarkup;
      modal.querySelector(".retail-markup").value =
        tableRow.dataset.retailMarkup;
      break;
    case "options":
      modal.querySelector(".id").value = tableRow.dataset.id;
      modal.querySelector(".name").value = tableRow.dataset.name;
      modal.querySelector(".wholesale-price").value =
        tableRow.dataset.wholesalePrice;
      modal.querySelector(".retail-price").value = tableRow.dataset.retailPrice;
      modal.querySelector(".image-path").value =
        tableRow.dataset.imagePath || "no_image";
      break;
    case "products":
      modal.querySelector(".id").value = tableRow.dataset.id;
      modal.querySelector(".name").value = tableRow.dataset.name;
      modal.querySelector(".wholesale-price").value =
        tableRow.dataset.wholesalePrice;
      modal.querySelector(".retail-price").value = tableRow.dataset.retailPrice;
      modal.querySelector(".image-path").value =
        tableRow.dataset.imagePath || "no_image";
      break;
    case "rails":
      modal.querySelector(".id").value = tableRow.dataset.id;
      modal.querySelector(".name").value = tableRow.dataset.name;
      modal.querySelector(".wholesale-price").value =
        tableRow.dataset.wholesalePrice;
      modal.querySelector(".retail-price").value = tableRow.dataset.retailPrice;
      modal.querySelector(".specifications").value =
        tableRow.dataset.specifications || "";
      modal.querySelector(".image-path").value =
        tableRow.dataset.imagePath || "no_image";
      break;
    case "residential_drives":
      modal.querySelector(".id").value = tableRow.dataset.id;
      modal.querySelector(".name").value = tableRow.dataset.name;
      modal.querySelector(".wholesale-price").value =
        tableRow.dataset.wholesalePrice;
      modal.querySelector(".retail-price").value = tableRow.dataset.retailPrice;
      modal.querySelector(".specifications").value =
        tableRow.dataset.specifications || "";
      modal.querySelector(".image-path").value =
        tableRow.dataset.imagePath || "no_image";
      break;
    case "users":
      modal.querySelector(".id").value = tableRow.dataset.id;
      modal.querySelector(".fullname").value = tableRow.dataset.fullname;
      modal.querySelector(".email").value = tableRow.dataset.email;
      modal.querySelector(".phone").value = tableRow.dataset.phone;
      modal.querySelector(".role").value = tableRow.dataset.role;
      modal.querySelector(".password").value = "";
      break;
    case "manual_drive_prices":
      modal.querySelector(".id").value = tableRow.dataset.id;
      modal.querySelector(".chain-meter-retail-price").value = tableRow.dataset.chainMeterRetailPrice;
      modal.querySelector(".chain-meter-wholesale-price").value = tableRow.dataset.chainMeterWholesalePrice;
      modal.querySelector(".rcp-retail-price").value = tableRow.dataset.rcpRetailPrice;
      modal.querySelector(".rcp-wholesale-price").value = tableRow.dataset.rcpWholesalePrice;
      break;
  }

  modal.showModal();
}

async function updateRow(event) {
  event.preventDefault();

  const form = event.target;
  const url = form.dataset.url;
  const formData = new FormData(form);
  const hasFileInput = form.querySelector('input[type="file"]') !== null;

  const response = await fetch(url, {
    method: "PUT",
    body: hasFileInput ? formData : new URLSearchParams(formData),
    headers: hasFileInput
      ? {}
      : { "Content-Type": "application/x-www-form-urlencoded" },
  });

  if (!response.ok) {
    const text = await response.text();
    document.querySelector('dialog[open]')?.close();
    Swal.fire({ icon: "error", title: "Ошибка", text });
    return;
  }

  window.location.reload();
}

async function deleteRow(rowId, tableName) {
  const result = await Swal.fire({
    title: "Удалить запись?",
    icon: "warning",
    showCancelButton: true,
    confirmButtonText: "Удалить",
    cancelButtonText: "Отмена",
  });
  if (!result.isConfirmed) return;

  const response = await fetch("/tables/" + tableName + "/" + rowId, {
    method: "DELETE",
  });

  if (!response.ok) {
    const text = await response.text();
    Swal.fire({ icon: "error", title: "Ошибка", text });
    return;
  }

  window.location.reload();
}
