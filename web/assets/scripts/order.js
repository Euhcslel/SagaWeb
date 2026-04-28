const Proto = {
  DocumentsList: null,
  UpdateOrderStatusRequest: null,
};

const OrderStatus = {
  pending: 0,
  confirmed: 1,
  paid: 2,
  done: 3,
  cancelled: 4,
};

// Функция, инициализирующая proto-схемы
async function initProtobuf() {
  const documentsRoot = await protobuf.load("/api/proto/documents.proto");
  Proto.DocumentsList = documentsRoot.lookupType("proto.DocumentsList");

  const statusRoot = await protobuf.load("/api/proto/status.proto");
  Proto.UpdateOrderStatusRequest = statusRoot.lookupType(
    "proto.UpdateOrderStatusRequest",
  );
}

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

function createHeaderRow(text) {
  const row = document.createElement("tr");
  const cell = document.createElement("td");
  cell.colSpan = 2;
  cell.textContent = text;
  row.appendChild(cell);
  return row;
}

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

window.downloadDocument = async (button) => {
  const documentName = button
    .closest("tr")
    .querySelector(".document-name").textContent;

  const documentType = button.closest("tbody").dataset.documentType;

  const response = await fetch(window.location.pathname + `/documents/${documentType}/${documentName}`, {
    method: "GET",
  });

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

window.deleteDocument = async (button) => {
  const documentName = button
    .closest("tr")
    .querySelector(".document-name").textContent;
  const documentType = button.closest("tbody").dataset.documentType;

  const res = await fetch(window.location.pathname + `/documents/${documentType}/${documentName}`, {
    method: "DELETE",
  });

  if (!res.ok) throw new Error("Delete request failed");

  window.location.reload();
};

await initProtobuf();
