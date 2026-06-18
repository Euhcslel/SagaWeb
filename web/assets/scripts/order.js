import {
  OrderStatus,
  removeItem,
  addProduct,
  updateProductsPrice,
  fmtPrice,
} from "./gate-form-utils.js";
import {
  create,
  toBinary,
  fromBinary,
  DocumentsListSchema,
  UpdateOrderStatusRequestSchema,
  UpdateProductsRequestSchema,
} from "./proto_bundle.js";

window.removeItem = removeItem;
window.addProduct = addProduct;
window.updateProductsPrice = updateProductsPrice;

const statusSelectEl = document.getElementById("status");
let currentOrderStatus = statusSelectEl ? statusSelectEl.value : null;

window.deleteGate = async (gateId) => {
  const result = await Swal.fire({
    title: "Удалить ворота?",
    text: "Вы уверены, что хотите удалить ворота из заказа?",
    icon: "question",
    showCancelButton: true,
    confirmButtonText: "Да, удалить",
    cancelButtonText: "Отмена",
  });

  if (result.isConfirmed) {
    try {
      const path = window.location.pathname;
      const parts = path.split("/");
      const saleId = parts[2];

      const res = await fetch("/orders/" + saleId + "/" + gateId, {
        method: "DELETE",
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(errorText || "Не удалось удалить ворота");
      }

      window.location.reload();
    } catch (error) {
      Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
    }
  }
};

window.deleteOrder = async (saleId) => {
  const result = await Swal.fire({
    title: "Удалить заказ?",
    text: "Вы уверены, что хотите удалить заказ?",
    icon: "question",
    showCancelButton: true,
    confirmButtonText: "Да, удалить",
    cancelButtonText: "Отмена",
  });

  if (result.isConfirmed) {
    try {
      const res = await fetch("/orders/" + saleId, {
        method: "DELETE",
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(errorText || "Не удалось удалить заказ");
      }

      window.location.href = "/orders";
    } catch (error) {
      Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
    }
  }
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
    } else {
      const errorText = await response.text();
      throw new Error(errorText || "Не удалось добавить ворота");
    }
  } catch (error) {
    Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
  }
};

window.loadDocuments = async () => {
  try {
    const res = await fetch(`${window.location.pathname}/documents`);

    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(errorText || "Не удалось загрузить документы");
    }

    const buf = await res.arrayBuffer();
    const docsList = fromBinary(DocumentsListSchema, new Uint8Array(buf));

    const docsTable = document.getElementById("documents-list");
    docsTable.innerHTML = "";

    const template = document.getElementById("document-info-template");

    if (docsList.offers != undefined && docsList.offers.length > 0) {
      const offerDiv = document.createElement("tbody");
      offerDiv.dataset.documentType = "offer";
      offerDiv.appendChild(createHeaderRow("Коммерческие предложения"));

      docsList.offers.forEach((offer) => {
        const clone = template.content.cloneNode(true);
        clone.querySelector(".document-name").textContent = offer.offerNumber;
        offerDiv.appendChild(clone);
      });

      docsTable.appendChild(offerDiv);
    }

    if (docsList.contracts != undefined && docsList.contracts.length > 0) {
      const contractDiv = document.createElement("tbody");
      contractDiv.dataset.documentType = "contract";
      contractDiv.appendChild(createHeaderRow("Договоры"));

      docsList.contracts.forEach((contract) => {
        const clone = template.content.cloneNode(true);
        clone.querySelector(".document-name").textContent = contract.contractNumber;
        contractDiv.appendChild(clone);
      });

      docsTable.appendChild(contractDiv);
    }

    if (docsList.bills != undefined && docsList.bills.length > 0) {
      const billDiv = document.createElement("tbody");
      billDiv.dataset.documentType = "bill";
      billDiv.appendChild(createHeaderRow("Счета"));

      docsList.bills.forEach((bill) => {
        const clone = template.content.cloneNode(true);
        clone.querySelector(".document-name").textContent = bill.billNumber;
        billDiv.appendChild(clone);
      });

      docsTable.appendChild(billDiv);
    }

    if (docsList.documents != undefined && docsList.documents.length > 0) {
      const documentDiv = document.createElement("tbody");
      documentDiv.dataset.documentType = "document";
      documentDiv.appendChild(createHeaderRow("Прочие документы"));

      docsList.documents.forEach((doc) => {
        const clone = template.content.cloneNode(true);
        clone.querySelector(".document-name").textContent = doc.name;
        documentDiv.appendChild(clone);
      });

      docsTable.appendChild(documentDiv);
    }

    if (docsTable.children.length === 0) {
      const row = document.createElement("tr");
      const cell = document.createElement("td");
      cell.textContent = "Документов нет";
      cell.style.padding = "1rem";
      cell.style.textAlign = "center";
      cell.style.opacity = "0.6";
      row.appendChild(cell);
      docsTable.appendChild(row);
    }

    document.getElementById("documents-modal").showModal();
  } catch (error) {
    Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
  }
};

function createHeaderRow(text) {
  const row = document.createElement("tr");
  const cell = document.createElement("td");
  cell.colSpan = 2;
  cell.textContent = text;
  row.appendChild(cell);
  return row;
}

window.updateUploadButton = () => {
  const hasFile =
    document.getElementById("upload-offer").files.length > 0 ||
    document.getElementById("upload-contract").files.length > 0 ||
    document.getElementById("upload-bill").files.length > 0 ||
    document.getElementById("upload-document").files.length > 0;
  document.getElementById("upload-docs-btn").disabled = !hasFile;
};

window.uploadDocuments = async () => {
  const orderId = window.location.pathname.split("/")[2];
  const uploads = [
    { inputId: "upload-offer", endpoint: `/orders/${orderId}/offer` },
    { inputId: "upload-contract", endpoint: `/orders/${orderId}/contract` },
    { inputId: "upload-bill", endpoint: `/orders/${orderId}/bill` },
    { inputId: "upload-document", endpoint: `/orders/${orderId}/document` },
  ];

  try {
    for (const { inputId, endpoint } of uploads) {
      const file = document.getElementById(inputId).files[0];
      if (!file) continue;

      const formData = new FormData();
      formData.append("file", file);

      const res = await fetch(endpoint, { method: "POST", body: formData });
      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(errorText || "Ошибка при загрузке файла");
      }
    }
    document.getElementById("add-documents-modal").close();
    window.location.reload();
  } catch (error) {
    document.querySelector('dialog[open]')?.close();
    Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
  }
};

const STATUS_REQUIRES_DOCS = new Set(["paid", "in_production", "ready", "done"]);

window.changeOrderStatus = async () => {
  const statusSelect = document.getElementById("status");
  const selectedStatus = statusSelect.options[statusSelect.selectedIndex].value;
  const selectedLabel = statusSelect.options[statusSelect.selectedIndex].text;

  const confirmed = await Swal.fire({
    title: "Изменить статус?",
    text: `Статус будет изменён на «${selectedLabel}»`,
    icon: "question",
    showCancelButton: true,
    confirmButtonText: "Да, изменить",
    cancelButtonText: "Отмена",
  });

  if (!confirmed.isConfirmed) {
    statusSelect.value = currentOrderStatus;
    return;
  }

  try {
    if (STATUS_REQUIRES_DOCS.has(selectedStatus)) {
      const orderId = window.location.pathname.split("/")[2];
      const docsRes = await fetch(`/orders/${orderId}/documents`);
      if (!docsRes.ok) throw new Error("Не удалось проверить документы заказа");

      const docsArr = await docsRes.arrayBuffer();
      const docsList = fromBinary(DocumentsListSchema, new Uint8Array(docsArr));

      if (!docsList.contracts || docsList.contracts.length === 0) {
        statusSelect.value = currentOrderStatus;
        Swal.fire({
          title: "Невозможно изменить статус",
          text: "К заказу не прикреплено приложение к договору",
          icon: "warning",
        });
        return;
      }

      if (!docsList.bills || docsList.bills.length === 0) {
        statusSelect.value = currentOrderStatus;
        Swal.fire({
          title: "Невозможно изменить статус",
          text: "К заказу не прикреплён счёт",
          icon: "warning",
        });
        return;
      }
    }

    const payload = { status: OrderStatus[selectedStatus] };
    const msg = create(UpdateOrderStatusRequestSchema, payload);
    const buf = toBinary(UpdateOrderStatusRequestSchema, msg);

    const res = await fetch(window.location.pathname + "/status", {
      method: "PUT",
      headers: { "Content-Type": "application/x-protobuf" },
      body: buf,
    });

    if (!res.ok) {
      statusSelect.value = currentOrderStatus;
      const errorText = await res.text();
      throw new Error(errorText || "Не удалось изменить статус");
    }

    window.location.reload();
  } catch (error) {
    statusSelect.value = currentOrderStatus;
    Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
  }
};

window.downloadDocument = async (button) => {
  try {
    const documentName = button
      .closest("tr")
      .querySelector(".document-name").textContent;

    const documentType = button.closest("tbody").dataset.documentType;

    const response = await fetch(
      window.location.pathname + `/documents/${documentType}/${documentName}`,
      { method: "GET" },
    );

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || "Не удалось скачать документ");
    }

    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = documentName;
    document.body.appendChild(a);
    a.click();
    a.remove();
    window.URL.revokeObjectURL(url);
  } catch (error) {
    document.querySelector('dialog[open]')?.close();
    Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
  }
};

window.deleteDocument = async (button) => {
  const documentsModal = button.closest("dialog");

  if (documentsModal) {
    documentsModal.close();
  }

  const result = await Swal.fire({
    title: "Удалить документ?",
    text: "Вы уверены, что хотите удалить этот документ?",
    icon: "question",
    showCancelButton: true,
    confirmButtonText: "Да, удалить",
    cancelButtonText: "Отмена",
  });

  if (!result.isConfirmed) {
    if (documentsModal) {
      documentsModal.showModal();
    }
    return;
  }

  try {
    const documentName = button
      .closest("tr")
      .querySelector(".document-name")
      .textContent.trim();

    const documentType = button.closest("tbody").dataset.documentType;

    const res = await fetch(
      window.location.pathname + `/documents/${documentType}/${documentName}`,
      { method: "DELETE" },
    );

    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(errorText || "Не удалось удалить документ");
    }

    window.location.reload();
  } catch (error) {
    Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
  }
};

window.saveProducts = async () => {
  try {
    const productList = document.getElementById("product-list");

    const productMap = new Map();

    productList.querySelectorAll(".product-item").forEach((item) => {
      const productId = Number(item.querySelector(".product").value);
      const amount = Number(item.querySelector(".amount").value);

      const currentAmount = productMap.get(productId) || 0;
      productMap.set(productId, currentAmount + amount);
    });

    const products = Array.from(productMap.entries()).map(
      ([productId, amount]) => ({
        productId: BigInt(productId),
        amount,
      }),
    );

    const payload = { products };

    const msg = create(UpdateProductsRequestSchema, payload);
    const buf = toBinary(UpdateProductsRequestSchema, msg);

    const path = window.location.pathname;

    const res = await fetch(path + "/products", {
      method: "PUT",
      headers: { "Content-Type": "application/x-protobuf" },
      body: buf,
    });

    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(errorText || "Не удалось сохранить товары");
    }

    window.location.reload();
  } catch (error) {
    Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
  }
};

window.updateOrderPrice = () => {
  let gatesRetailPrice = 0;
  let gatesWholesalePrice = 0;

  const gates = document.getElementsByClassName("gate");
  [...gates].forEach((gate) => {
    const gateWholesalePrice = Number(gate.dataset.wholesalePrice);
    const gateRetailPrice = Number(gate.dataset.retailPrice);
    const amount = Number(gate.dataset.amount);

    gatesRetailPrice += gateRetailPrice * amount;
    gatesWholesalePrice += gateWholesalePrice * amount;
  });

  const productList = document.getElementById("product-list");
  const productsRetailPrice = Number(productList.dataset.retailPrice);
  const productsWholesalePrice = Number(productList.dataset.wholesalePrice);

  const orderWholesalePriceElement = document.getElementById("order-wholesale-price");
  orderWholesalePriceElement.textContent = fmtPrice(
    productsWholesalePrice + gatesWholesalePrice,
  );

  const orderRetailPriceElement = document.getElementById("order-retail-price");
  orderRetailPriceElement.textContent = fmtPrice(
    productsRetailPrice + gatesRetailPrice,
  );
};

window.downloadAppendice = async () => {
  try {
    const orderId = window.location.pathname.split("/")[2];

    const response = await fetch(`/orders/${orderId}/appendice`, { method: "GET" });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || "Ошибка при генерации приложения");
    }

    const blob = await response.blob();
    const disposition = response.headers.get("Content-Disposition") || "";
    const match = disposition.match(/filename="([^"]+)"/);
    const filename = match ? match[1] : `appendice_${orderId}.pdf`;

    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    a.remove();
    window.URL.revokeObjectURL(url);

    document.getElementById("generate-modal").close();
  } catch (error) {
    document.querySelector('dialog[open]')?.close();
    Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
  }
};

window.openDocumentModal = (type) => {
  document.getElementById("generate-modal").close();
  document.getElementById("offer-modal").showModal();
};

window.downloadOffer = async (type) => {
  const clientName = document.getElementById("offer-client-name").value.trim();
  if (!clientName) {
    Swal.fire({ title: "Ошибка", text: "Укажите ФИО покупателя", icon: "error" });
    return;
  }

  try {
    const orderId = window.location.pathname.split("/")[2];
    const params = new URLSearchParams({ client_name: clientName });
    const response = await fetch(`/orders/${orderId}/offer/${type}?${params}`, { method: "GET" });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || "Ошибка при генерации КП");
    }

    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `offer_${orderId}_${type}.pdf`;
    document.body.appendChild(a);
    a.click();
    a.remove();
    window.URL.revokeObjectURL(url);

    document.getElementById("offer-modal").close();
  } catch (error) {
    document.querySelector('dialog[open]')?.close();
    Swal.fire({ title: "Ошибка", text: error.message, icon: "error" });
  }
};

updateProductsPrice();
updateOrderPrice();
