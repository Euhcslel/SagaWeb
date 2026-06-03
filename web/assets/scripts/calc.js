import {
  gateTypeProto,
  removeItem,
  addGateOption,
  addProduct,
  updateOptionsPrice,
  updateProductsPrice,
  updateGatePrice as sharedUpdateGatePrice,
  fmtPrice,
} from "./gate-form-utils.js";

window.removeItem = removeItem;
window.addGateOption = addGateOption;
window.updateOptionsPrice = updateOptionsPrice;
window.addProduct = addProduct;
window.updateProductsPrice = updateProductsPrice;
window.updateGatePrice = (gate) => {
  sharedUpdateGatePrice(gate);
  updateOrderPrice();
};

// Proto-схемы
let Proto = {
  SizePrice: null,
  OrderRequest: null,
};

// Функция, инициализирующая proto-схемы
async function initProtobuf() {
  const pricesRoot = await protobuf.load("/api/proto/prices.proto");
  Proto.SizePrice = pricesRoot.lookupType("proto.SizePrice");

  const orderRoot = await protobuf.load("/api/proto/order.proto");
  Proto.OrderRequest = orderRoot.lookupType("proto.OrderRequest");
}

// Функция, удаляющая вкладку ворот и соответствующий шаблон формы
window.removeGateElement = async (event) => {
  const tabElement = event.currentTarget.closest("#gate-tabs > li");

  const result = await Swal.fire({
    title: "Удалить ворота?",
    text: "Вы уверены, что хотите удалить ворота из заказа?",
    icon: "question",
    showCancelButton: true,
    confirmButtonText: "Да, удалить",
    cancelButtonText: "Отмена",
  });

  if (result.isConfirmed) {
    const gateTabs = document.getElementById("gate-tabs");
    const gateList = document.getElementById("gate-list");

    const tabIndex = Array.from(gateTabs.children).indexOf(tabElement);

    tabElement.remove();
    gateList.children[tabIndex]?.remove();

    const firstInput = gateTabs.querySelector(
      'li:first-child input[type="radio"]',
    );
    if (firstInput) {
      firstInput.checked = true;
    }

    updateOrderPrice();
  }
};

// Функция, вызывающаяся при смене текущих ворот
// необходима, чтобы корректно менять цену ворот в заказе
window.onChangeGateButton = (e) => {
  if (!e.target.matches('input[name="gates"]')) return;

  const radios = document.querySelectorAll('#gate-tabs input[name="gates"]');
  const index = [...radios].findIndex((r) => r.checked);

  const gate = document.getElementsByClassName("gate-item")[index];

  document.getElementById("gate-retail-price").textContent = fmtPrice(
    Number(gate?.dataset.retailPrice || 0),
  );

  document.getElementById("gate-wholesale-price").textContent = fmtPrice(
    Number(gate?.dataset.wholesalePrice || 0),
  );
};

// Функция, добавляющая новую вкладку и шаблон формы
window.addGateElement = () => {
  const gateItems = document.getElementsByClassName("gate-item");
  if (gateItems.length === 15) {
    Swal.fire("Нельзя добавить больше 15 ворот в заказ.");
    return;
  }

  const gateTemplate = document.getElementById("gate-template");
  const gateList = document.getElementById("gate-list");

  const gateClone = gateTemplate.content.cloneNode(true);
  gateList.append(gateClone);

  const gateTabElementTemplate = document.getElementById("gate-tab-template");
  const gateTabs = document.getElementById("gate-tabs");

  const gateTabElementClone = gateTabElementTemplate.content.cloneNode(true);
  gateTabs.append(gateTabElementClone);

  const lastGate = gateList.lastElementChild;
  updateGateType(lastGate.querySelector(".gate-type"));
  updateGateSizePrice(lastGate.querySelector(".width"));
  updateDriveType(lastGate.querySelector(".drive-type"));
};

// Функция, которая срабатывает при смене типа ворот
// обновляет динамические значения формы, которые зависят от типа ворот
window.updateGateType = async (select) => {
  const gate = select.closest(".gate-item");
  const cfg = configuration[select.value + "Configuration"];

  const width = gate.querySelector(".width");
  width.min = cfg.widthParams.MinValue;
  width.max = cfg.widthParams.MaxValue;
  width.value = cfg.widthParams.MinValue;

  updateGatePreview(gate);

  const height = gate.querySelector(".height");
  height.min = cfg.heightParams.MinValue;
  height.max = cfg.heightParams.MaxValue;
  height.value = cfg.heightParams.MinValue;

  const liftType = gate.querySelector(".lift-type");
  fillSelect(liftType, cfg.liftTypes, "markup");

  const cycleAmount = gate.querySelector(".cycle-amount");
  fillSelect(cycleAmount, cfg.cycleAmounts, "markup", "Amount");

  const driveType = gate.querySelector(".drive-type");
  fillSelect(driveType, cfg.driveTypes, "object");

  const drive = gate.querySelector(".drive");
  fillSelect(drive, cfg.drives, "price");
  drive.selectedIndex = 0;

  if (select.value === "residential") {
    const rail = gate.querySelector(".rail");
    fillSelect(rail, cfg.rails, "price");
    rail.selectedIndex = 0;
  }

  updateDriveType(gate.querySelector(".drive-type"));
  await updateGateSizePrice(width);
};

// Функция, которая используется для заполнения элементов select
// в зависимости от переданного типа, создает option с нужными параметрами
function fillSelect(select, items, type, labelField = "Name") {
  select.innerHTML = "";
  if (type === "markup") {
    items.forEach((item) => {
      const option = document.createElement("option");
      option.dataset.retailMarkup = item.RetailMarkup;
      option.dataset.wholesaleMarkup = item.WholesaleMarkup;
      option.value = item.ID;
      option.textContent = item[labelField];

      select.append(option);
    });
  } else if (type === "object") {
    Object.entries(items).forEach(([key, value]) => {
      const option = document.createElement("option");
      option.value = key;
      option.textContent = value;

      select.append(option);
    });
  } else if (type === "price") {
    items.forEach((item) => {
      const option = document.createElement("option");
      option.dataset.retailPrice = item.RetailPrice;
      option.dataset.wholesalePrice = item.WholesalePrice;
      option.value = item.ID;
      option.textContent = item.Name;

      select.append(option);
    });
  }
}

// Функция, срабатывающая при изменении привода
// показывает нужные поля ввода и скрывает ненужные
window.updateDriveType = (select) => {
  const gate = select.closest(".gate-item");

  const isManual = select.value === "manual";

  const chainLabel = gate.querySelector(".chain-length").closest("label");
  const driveAutoBlock = gate.querySelector(".drive-auto");
  const railLabel = gate.querySelector(".rail").closest("label");
  const gateType = gate.querySelector(".gate-type").value;

  chainLabel.hidden = !isManual;
  driveAutoBlock.hidden = isManual;
  railLabel.hidden = isManual || gateType !== "residential";

  updateGatePrice(gate);
};

// Функция, срабатывающая при изменении значений ширины и высоты
// запрашивает цены клиента и дилера у сервера и присваивает их блоку sizes
window.updateGateSizePrice = async (input) => {
  const gate = input.closest(".gate-item");
  const width = gate.querySelector(".width");
  const height = gate.querySelector(".height");
  const gateType = gate.querySelector(".gate-type");

  const sizes = input.closest(".sizes");

  if (
    Number(width.value) < Number(width.min) ||
    Number(width.value) > Number(width.max) ||
    Number(height.value) < Number(height.min) ||
    Number(height.value) > Number(height.max)
  ) {
    return;
  }

  updateGatePreview(gate);

  const price = await fetchSizePrice(width.value, height.value, gateType.value);

  sizes.dataset.retailPrice = price.retail;
  sizes.dataset.wholesalePrice = price.wholesale;

  updateGatePrice(gate);
};

function getSelectedGateColor(gate) {
  const colorSelect = gate.querySelector(".color-out");
  const selectedColor = colorSelect?.options[colorSelect.selectedIndex];

  const hex = selectedColor?.dataset.hex;

  if (!hex) {
    return "#f0efea";
  }

  return hex.startsWith("#") ? hex : `#${hex}`;
}

function drawGateLines(linesGroup, x, y, width, height) {
  linesGroup.replaceChildren();

  // Чем выше ворота, тем больше секций, но без перебора.
  const sections = Math.max(3, Math.min(7, Math.round(height / 32)));

  for (let i = 1; i < sections; i++) {
    const lineY = y + (height / sections) * i;

    const line = document.createElementNS("http://www.w3.org/2000/svg", "line");

    line.setAttribute("x1", x + 4);
    line.setAttribute("y1", lineY);
    line.setAttribute("x2", x + width - 4);
    line.setAttribute("y2", lineY);

    linesGroup.append(line);
  }
}

function updateGatePreview(gate) {
  const width = Number(gate.querySelector(".width").value || 2000);
  const height = Number(gate.querySelector(".height").value || 1800);

  const preview = gate.querySelector(".gate-preview");
  const rect = gate.querySelector(".gate-rect");
  const linesGroup = gate.querySelector(".gate-lines");
  const widthLabel = gate.querySelector(".width-label");
  const heightLabel = gate.querySelector(".height-label");

  preview.style.setProperty("--gate-color", getSelectedGateColor(gate));

  // Рабочая зона ворот. Справа оставляем место под H, снизу под B.
  const areaX = 35;
  const areaY = 45;
  const areaW = 330;
  const areaH = 190;

  const ratio = width / height;

  let previewW;
  let previewH;

  if (ratio > areaW / areaH) {
    previewW = areaW;
    previewH = areaW / ratio;
  } else {
    previewH = areaH;
    previewW = areaH * ratio;
  }

  const x = areaX + (areaW - previewW) / 2;
  const y = areaY + (areaH - previewH) / 2;

  rect.setAttribute("x", x);
  rect.setAttribute("y", y);
  rect.setAttribute("width", previewW);
  rect.setAttribute("height", previewH);

  drawGateLines(linesGroup, x, y, previewW, previewH);

  widthLabel.setAttribute("x", x + previewW / 2);
  widthLabel.setAttribute("y", y + previewH + 34);
  widthLabel.setAttribute("text-anchor", "middle");
  widthLabel.textContent = `B = ${width}`;

  heightLabel.setAttribute("x", 435);
  heightLabel.setAttribute("y", y + previewH / 2);
  heightLabel.setAttribute("text-anchor", "middle");
  heightLabel.setAttribute("dominant-baseline", "middle");
  heightLabel.textContent = `H = ${height}`;
}

window.updateGatePreview = updateGatePreview;


// Функция для скачивания коммерческого предложения
// Собирает все данные с формы, кодирует для protobuf и отправляет на сервер
// Получает pdf-файл и скачивает его
window.downloadOffer = async () => {
  const payload = {
    orderGates: buildGatePayload(),
    products: buildProductsList(),
  };

  const err = Proto.OrderRequest.verify(payload);
  if (err) throw new Error(err);

  const msg = Proto.OrderRequest.create(payload);
  const buf = Proto.OrderRequest.encode(msg).finish();

  const res = await fetch("/offer", {
    method: "POST",
    headers: { "Content-Type": "application/x-protobuf" },
    body: buf,
  });

  if (!res.ok) {
    throw new Error("Не удалось сформировать КП");
  }

  const blob = await res.blob();

  const url = URL.createObjectURL(blob);

  const a = document.createElement("a");
  a.href = url;
  a.download = "commercial-offer.pdf";
  a.click();

  URL.revokeObjectURL(url);
};

// Функция, запрашивающая цену за размер ворот
// принимает параметры: w - ширина, h - высота, t - тип ворот
async function fetchSizePrice(w, h, t) {
  const res = await fetch(`/sizes?width=${w}&height=${h}&gateType=${t}`);

  const buf = await res.arrayBuffer();

  const msg = Proto.SizePrice.decode(new Uint8Array(buf));

  const d = Proto.SizePrice.toObject(msg, { longs: Number });

  if (d.dealer) {
    return {
      retail: d.dealer.clientPrice / 100,
      wholesale: d.dealer.dealerPrice / 100,
    };
  }

  return { retail: d.client.clientPrice / 100, wholesale: 0 };
}

// Функция, которая пересчитывает стоимость всего заказа
window.updateOrderPrice = () => {
  let totalRetailPrice = 0;
  let totalWholesalePrice = 0;

  const gates = document.getElementsByClassName("gate-item");
  [...gates].forEach((gate) => {
    const amount = Number(gate.querySelector(".gate-amount").value);

    const gateRetailPrice = Number(gate.dataset.retailPrice || 0);
    totalRetailPrice += gateRetailPrice * amount;

    const gateWholesalePrice = Number(gate.dataset.wholesalePrice || 0);
    totalWholesalePrice += gateWholesalePrice * amount;
  });

  const products = document.getElementById("product-list");
  totalRetailPrice += Number(products.dataset.retailPrice || 0);
  totalWholesalePrice += Number(products.dataset.wholesalePrice || 0);

  const orderRetailPriceElement = document.getElementById("order-retail-price");
  orderRetailPriceElement.textContent = fmtPrice(totalRetailPrice);

  const orderWholesalePriceElement = document.getElementById(
    "order-wholesale-price",
  );
  orderWholesalePriceElement.textContent = fmtPrice(totalWholesalePrice);
};

// Функция, формирующая список товаров для отправки заказа
function buildProductsList() {
  const productsMap = {};

  const productItems = document.getElementsByClassName("product-item");

  [...productItems].forEach((productItem) => {
    const productSelect = productItem.querySelector(".product");
    const selectedProduct = productSelect.options[productSelect.selectedIndex];

    const productId = Number(selectedProduct.value);
    const amount = Number(productItem.querySelector(".amount").value);

    if (productsMap[productId] === undefined) {
      productsMap[productId] = 0;
    }

    productsMap[productId] += amount;
  });

  return Object.entries(productsMap).map(([productId, amount]) => ({
    productId: Number(productId),
    amount: amount,
  }));
}

// Функция, пракильно формирующая описание типа привода и информацию о нем.
function buildDrivePayload(gate) {
  const driveType = getSelectedOption(gate.querySelector(".drive-type")).value;
  if (driveType === "manual") {
    return {
      manual: {
        chainLength: Number(gate.querySelector(".chain-length").value),
      },
    };
  }

  if (driveType === "industrial") {
    return {
      industrial: {
        driveId: Number(gate.querySelector(".drive").value),
      },
    };
  }

  if (driveType === "residential") {
    return {
      residential: {
        driveId: Number(gate.querySelector(".drive").value),
        railId: Number(gate.querySelector(".rail").value),
      },
    };
  }

  return null;
}

// Функция, возвращающа выбранный option у select
function getSelectedOption(select) {
  return select.options[select.selectedIndex];
}

// Функция для формирования списка ворот с заказе.
// Возвращает массив объектов ворот
function buildGatePayload() {
  const orderGates = [];
  const gates = document.getElementsByClassName("gate-item");
  [...gates].forEach((gate) => {
    // Тип ворот
    const gateType =
      gateTypeProto[getSelectedOption(gate.querySelector(".gate-type")).value];

    // Ширина и высота
    const width = Number(gate.querySelector(".width").value);
    const height = Number(gate.querySelector(".height").value);

    // Высота притолки
    const headroom = Number(gate.querySelector(".headroom").value);

    // Тип подъема
    const liftType = Number(
      getSelectedOption(gate.querySelector(".lift-type")).value,
    );

    // Количество циклов
    const cycleAmount = Number(
      getSelectedOption(gate.querySelector(".cycle-amount")).value,
    );

    // Цвет снаружи
    const colorOut = Number(
      getSelectedOption(gate.querySelector(".color-out")).value,
    );

    // Привод
    const drive = buildDrivePayload(gate);

    // Количество ворот
    const gateAmount = Number(gate.querySelector(".gate-amount").value);

    // Сбор дополнительных опций
    const optionsMap = {};

    const optionItems = gate.getElementsByClassName("option-item");

    [...optionItems].forEach((optionItem) => {
      const optionSelect = optionItem.querySelector(".option");
      const selectedOption = optionSelect.options[optionSelect.selectedIndex];
      const optionId = Number(selectedOption.value);
      const amount = Number(optionItem.querySelector(".amount").value);

      if (optionsMap[optionId] === undefined) {
        optionsMap[optionId] = 0;
      }

      optionsMap[optionId] += amount;
    });

    const options = Object.entries(optionsMap).map(([optionId, amount]) => ({
      optionId: Number(optionId),
      amount: amount,
    }));

    orderGates.push({
      gateType: gateType,
      width: width,
      height: height,
      liftTypeId: liftType,
      colorOutId: colorOut,
      cycleAmountId: cycleAmount,
      options: options,
      headroom: headroom,
      drive: drive,
      amount: gateAmount,
    });
  });

  return orderGates;
}

// Функция дя оформления заказа.
// Собирает все данные с формы, кодирует для protobuf и отправляет на сервер
window.placeOrder = async () => {
  const result = await Swal.fire({
    title: "Оформить заказ?",
    text: "Вы готовы оформить заказ?",
    icon: "question",
    showCancelButton: true,
    confirmButtonText: "Да, оформить заказ",
    cancelButtonText: "Отмена",
  });

  if (result.isConfirmed) {
    const payload = {
      orderGates: buildGatePayload(),
      products: buildProductsList(),
    };

    const err = Proto.OrderRequest.verify(payload);
    if (err) throw new Error(err);

    const msg = Proto.OrderRequest.create(payload);
    const buf = Proto.OrderRequest.encode(msg).finish();

    const res = await fetch("/orders", {
      method: "POST",
      headers: { "Content-Type": "application/x-protobuf" },
      body: buf,
    });

    if (!res.ok) throw new Error("Order request failed");

    window.location.href = "/orders";
  }
};

// Инициализация proto-схем после загрузки
await initProtobuf();
// Чтобы изначально были хотя бы одни ворота
addGateElement();
