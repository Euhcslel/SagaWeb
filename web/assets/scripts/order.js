import {
  OrderStatus,
  removeItem,
  addProduct,
  updateProductsPrice,
} from "./gate-form-utils.js";

window.removeItem = removeItem;
window.addProduct = addProduct;
window.updateProductsPrice = updateProductsPrice;

const Proto = {
  DocumentsList: null,
  UpdateOrderStatusRequest: null,
  UpdateProductsRequest: null,
};

// Функция, инициализирующая proto-схемы
async function initProtobuf() {
  const documentsRoot = await protobuf.load("/api/proto/documents.proto");
  Proto.DocumentsList = documentsRoot.lookupType("proto.DocumentsList");

  const statusRoot = await protobuf.load("/api/proto/status.proto");
  Proto.UpdateOrderStatusRequest = statusRoot.lookupType(
    "proto.UpdateOrderStatusRequest",
  );

  const productsRoot = await protobuf.load("/api/proto/order.proto");
  Proto.UpdateProductsRequest = productsRoot.lookupType(
    "proto.UpdateProductsRequest",
  );
}

// Функция, удаляющая ворота из заказа
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

// Функция, удаляющая заказ пользователя
window.deleteOrder = async (saleId) => {
  const res = await fetch("/orders/" + saleId, {
    method: "DELETE",
  });

  if (!res.ok) throw new Error("Delete request failed");

  window.location.href = "/orders";
};

// Функция, добавляющая новые ворота в заказ
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

// Функция, подгружающая список документов
window.loadDocuments = async () => {
  const res = await fetch(`${window.location.pathname}/documents`);

  const buf = await res.arrayBuffer();

  const msg = Proto.DocumentsList.decode(new Uint8Array(buf));

  const docsList = Proto.DocumentsList.toObject(msg, { longs: Number });

  const docsTable = document.getElementById("documents-list");
  docsTable.innerHTML = "";

  const template = document.getElementById("document-info-template");
  let clone = template.content.cloneNode(true);

  if (docsList.offers != undefined) {
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

  if (docsList.contracts != undefined) {
    const contractDiv = document.createElement("tbody");
    contractDiv.dataset.documentType = "contract";
    contractDiv.appendChild(createHeaderRow("Договоры"));

    docsList.contracts.forEach((contract) => {
      const clone = template.content.cloneNode(true);
      clone.querySelector(".document-name").textContent =
        contract.contractNumber;

      contractDiv.appendChild(clone);
    });

    docsTable.appendChild(contractDiv);
  }

  if (docsList.bills != undefined) {
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

  if (docsList.documents != undefined) {
    const documentDiv = document.createElement("tbody");
    documentDiv.dataset.documentType = "document";
    documentDiv.appendChild(createHeaderRow("Прочие документы"));

    docsList.documents.forEach((document) => {
      const clone = template.content.cloneNode(true);
      clone.querySelector(".document-name").textContent =
        document.documentNumber;

      documentDiv.appendChild(clone);
    });

    docsTable.appendChild(documentDiv);
  }

  document.getElementById("documents-modal").showModal();
};

// Функция, создающая строку заголовка для таблицы документов
function createHeaderRow(text) {
  const row = document.createElement("tr");
  const cell = document.createElement("td");
  cell.colSpan = 2;
  cell.textContent = text;
  row.appendChild(cell);
  return row;
}

// Функция, изменяющая статус заказа
window.changeOrderStatus = async () => {
  const statusSelect = document.getElementById("status");
  const selectedStatus = statusSelect.options[statusSelect.selectedIndex].value;
  const payload = {
    status: OrderStatus[selectedStatus],
  };
  // Кодирование в protobuf
  const err = Proto.UpdateOrderStatusRequest.verify(payload);
  if (err) throw new Error(err);

  const msg = Proto.UpdateOrderStatusRequest.create(payload);
  const buf = Proto.UpdateOrderStatusRequest.encode(msg).finish();

  const path = window.location.pathname;

  const res = await fetch(path + "/status", {
    method: "PUT",
    headers: { "Content-Type": "application/x-protobuf" },
    body: buf,
  });

  if (!res.ok) throw new Error("Order request failed");

  window.location.reload();
};

// Функция, скачивающая документ
window.downloadDocument = async (button) => {
  const documentName = button
    .closest("tr")
    .querySelector(".document-name").textContent;

  const documentType = button.closest("tbody").dataset.documentType;

  const response = await fetch(
    window.location.pathname + `/documents/${documentType}/${documentName}`,
    {
      method: "GET",
    },
  );

  if (!response.ok) throw new Error("Download request failed");

  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = documentName;
  document.body.appendChild(a);
  a.click();
  a.remove();
  window.URL.revokeObjectURL(url);
};

// Функция, удаляющая документ
window.deleteDocument = async (button) => {
  const documentName = button
    .closest("tr")
    .querySelector(".document-name").textContent;
  const documentType = button.closest("tbody").dataset.documentType;

  const res = await fetch(
    window.location.pathname + `/documents/${documentType}/${documentName}`,
    {
      method: "DELETE",
    },
  );

  if (!res.ok) throw new Error("Delete request failed");

  window.location.reload();
};

// Функция, сохраниющая изменения товаров в заказе
window.saveProducts = async () => {
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
      productId,
      amount,
    }),
  );

  const payload = {
    products,
  };

  const err = Proto.UpdateProductsRequest.verify(payload);
  if (err) {
    throw new Error(err);
  }

  const msg = Proto.UpdateProductsRequest.create(payload);
  const buf = Proto.UpdateProductsRequest.encode(msg).finish();

  const path = window.location.pathname;

  const res = await fetch(path + "/products", {
    method: "PUT",
    headers: {
      "Content-Type": "application/x-protobuf",
    },
    body: buf,
  });

  if (!res.ok) {
    throw new Error("Order request failed");
  }

  window.location.reload();
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

  console.log(gatesWholesalePrice, gatesRetailPrice);
  const productList = document.getElementById("product-list");
  const productsRetailPrice = Number(productList.dataset.retailPrice);
  const productsWholesalePrice = Number(productList.dataset.wholesalePrice);
  console.log(productsRetailPrice, productsWholesalePrice);

  const orderWholesalePriceElement = document.getElementById(
    "order-wholesale-price",
  );
  orderWholesalePriceElement.textContent =
    productsWholesalePrice + gatesWholesalePrice;

  const orderRetailPriceElement = document.getElementById("order-retail-price");
  orderRetailPriceElement.textContent = productsRetailPrice + gatesRetailPrice;
};

await initProtobuf();
updateOrderPrice();
