/*
================================
GLOBAL STATE
================================
*/

const state = {
  role: null,
  gateType: null,
  products: [],
  orderPrice: { client: 0, dealer: 0 },
};

const orderGates = [];
let currentGateIndex = 0;

let Proto = {
  SizePrice: null,
  OrderRequest: null,
};

const GateType = {
  ind: 0,
  res: 1,
};

let DOM = {};

/*
================================
DOM CACHE
================================
*/

/* Кэширует все используемые DOM-элементы */
function initDomRefs() {
  DOM = {
    isDealer: document.getElementById("isDealer"),

    width: document.getElementById("width"),
    height: document.getElementById("height"),
    headroom: document.getElementById("headroom"),
    liftType: document.getElementById("liftType"),
    cycleAmount: document.getElementById("cycleAmount"),
    drive: document.getElementById("drive"),
    colorOut: document.getElementById("colorOut"),

    optionList: document.getElementById("optionList"),
    productList: document.getElementById("productList"),
    gateList: document.getElementById("gateList"),

    optionTemplate: document.getElementById("optionTemplate"),
    productTemplate: document.getElementById("productTemplate"),

    addGateBtn: document.getElementById("addGateToOrderButton"),
    addOptionBtn: document.getElementById("addOptionButton"),
    addProductBtn: document.getElementById("addProductButton"),

    gateClientPrice: document.querySelector(".gateClientPrice"),
    gateDealerPrice: document.querySelector(".gateDealerPrice"),
    totalClientPrice: document.querySelector(".totalClientPrice"),
    totalDealerPrice: document.querySelector(".totalDealerPrice"),

    headerGateNumber: document.querySelector(".headerGateNumber"),
  };
}

/*
================================
INITIALIZATION
================================
*/

/* Загружает protobuf-схемы для API */
async function initProtobuf() {
  const pricesRoot = await protobuf.load("/assets/proto_files/prices.proto");
  Proto.SizePrice = pricesRoot.lookupType("proto.SizePrice");

  const orderRoot = await protobuf.load("/assets/proto_files/order.proto");
  Proto.OrderRequest = orderRoot.lookupType("proto.OrderRequest");
}

/* Инициализирует калькулятор */
async function initCalculator() {
  initDomRefs();
  initGateTypeFromUrl();
  initUserRole();
  initEventListeners();

  orderGates.push(createGateFromDefaults());
  currentGateIndex = 0;

  await updateCurrentGateSizePrice();
  recalculatePrices();
  renderCalculator();
}

/* Определяет роль пользователя (dealer/client) */
function initUserRole() {
  const isDealer = DOM.isDealer?.value === "true";
  state.role = isDealer ? "dealer" : "client";
}

/* Читает тип ворот из URL */
function initGateTypeFromUrl() {
  const params = new URLSearchParams(window.location.search);
  state.gateType = params.get("gateType");
}

/*
================================
EVENTS
================================
*/

/* Регистрирует обработчики событий формы */
function initEventListeners() {
  const debouncedSizeUpdate = debounce(updateSize, 300);

  DOM.width.addEventListener("input", debouncedSizeUpdate);
  DOM.height.addEventListener("input", debouncedSizeUpdate);

  DOM.liftType.addEventListener("change", updateForm);
  DOM.cycleAmount.addEventListener("change", updateForm);
  DOM.drive.addEventListener("change", updateForm);
  DOM.colorOut.addEventListener("change", updateForm);

  DOM.optionList.addEventListener("input", updateForm);
  DOM.productList.addEventListener("input", updateProducts);

  DOM.addGateBtn.addEventListener("click", addGate);
  DOM.addOptionBtn.addEventListener("click", addOptionRow);
  DOM.addProductBtn.addEventListener("click", addProductRow);
}

/*
================================
CURRENT GATE
================================
*/

/* Возвращает текущие ворота из массива */
function currentGate() {
  return orderGates[currentGateIndex];
}

/*
================================
DEFAULT GATE
================================
*/

/* Значение ширины по умолчанию */
function getDefaultWidth() {
  return Number(DOM.width.min || DOM.width.value || 0);
}

/* Значение высоты по умолчанию */
function getDefaultHeight() {
  return Number(DOM.height.min || DOM.height.value || 0);
}

/* Значение headroom по умолчанию */
function getDefaultHeadroom() {
  return Number(DOM.headroom.min || DOM.headroom.value || 0);
}

/* Возвращает первый option select */
function getFirstOption(selectEl) {
  return selectEl.options[0] || null;
}

/* Создаёт объект ворот с начальными значениями */
function createGateFromDefaults() {
  const lift = getFirstOption(DOM.liftType);
  const cycle = getFirstOption(DOM.cycleAmount);
  const drive = getFirstOption(DOM.drive);
  const color = getFirstOption(DOM.colorOut);

  return {
    gateType: state.gateType,

    size: {
      width: getDefaultWidth(),
      height: getDefaultHeight(),
    },

    headroom: getDefaultHeadroom(),

    sizePrice: { client: 0, dealer: 0 },

    liftType: Number(lift?.value || 0),
    liftMarkup: {
      client: Number(lift?.dataset.clientMarkup || 0),
      dealer: Number(lift?.dataset.dealerMarkup || 0),
    },

    cycleAmount: Number(cycle?.value || 0),
    cycleMarkup: {
      client: Number(cycle?.dataset.clientMarkup || 0),
      dealer: Number(cycle?.dataset.dealerMarkup || 0),
    },

    driveId: Number(drive?.value || 0),
    drivePrice: {
      client: Number(drive?.dataset.clientPrice || 0),
      dealer: Number(drive?.dataset.dealerPrice || 0),
    },

    colorOutId: Number(color?.value || 0),

    options: [],
    gatePrice: { client: 0, dealer: 0 },
  };
}

/*
================================
FORM PARSING
================================
*/

/* Читает значения формы и обновляет текущие ворота */
function parseCurrentGateFromForm() {
  const gate = currentGate();

  gate.size.width = Number(DOM.width.value);
  gate.size.height = Number(DOM.height.value);
  gate.headroom = Number(DOM.headroom.value);

  const lift = DOM.liftType.selectedOptions[0];
  gate.liftType = Number(lift.value);
  gate.liftMarkup = {
    client: Number(lift.dataset.clientMarkup || 0),
    dealer: Number(lift.dataset.dealerMarkup || 0),
  };

  const cycle = DOM.cycleAmount.selectedOptions[0];
  gate.cycleAmount = Number(cycle.value);
  gate.cycleMarkup = {
    client: Number(cycle.dataset.clientMarkup || 0),
    dealer: Number(cycle.dataset.dealerMarkup || 0),
  };

  const drive = DOM.drive.selectedOptions[0];
  gate.driveId = Number(drive.value);
  gate.drivePrice = {
    client: Number(drive.dataset.clientPrice || 0),
    dealer: Number(drive.dataset.dealerPrice || 0),
  };

  gate.colorOutId = Number(DOM.colorOut.value);

  gate.options = [];

  DOM.optionList.querySelectorAll(".option-div").forEach((div) => {
    const select = div.querySelector("select");
    const input = div.querySelector("input");
    const option = select.options[select.selectedIndex];

    gate.options.push({
      id: Number(option.value),
      amount: getAmountFromInput(input),
      price: {
        client: Number(option.dataset.clientPrice || 0),
        dealer: Number(option.dataset.dealerPrice || 0),
      },
    });
  });
}

/* Читает список товаров из формы */
function readProductsFromForm() {
  state.products = [];

  DOM.productList.querySelectorAll(".product-div").forEach((div) => {
    const select = div.querySelector("select");
    const input = div.querySelector("input");
    const option = select.options[select.selectedIndex];

    state.products.push({
      id: Number(option.value),
      amount: getAmountFromInput(input),
      price: {
        client: Number(option.dataset.clientPrice || 0),
        dealer: Number(option.dataset.dealerPrice || 0),
      },
    });
  });
}

/*
================================
UPDATE
================================
*/

/* Обновляет цену при изменении размеров */
async function updateSize() {
  parseCurrentGateFromForm();
  await updateCurrentGateSizePrice();
  recalculatePrices();
}

/* Обновляет данные ворот */
function updateForm() {
  parseCurrentGateFromForm();
  recalculatePrices();
}

/* Обновляет список товаров */
function updateProducts() {
  readProductsFromForm();
  recalculatePrices();
}

/* Получает цену ворот по размеру */
async function updateCurrentGateSizePrice() {
  const gate = currentGate();
  gate.sizePrice = await fetchSizePrice(gate.size.width, gate.size.height);
}

/*
================================
API
================================
*/

/* Запрашивает цену размера ворот у сервера */
async function fetchSizePrice(width, height) {
  const res = await fetch(
    `/sizes?width=${width}&height=${height}&gateType=${state.gateType}`,
  );

  if (!res.ok) {
    throw new Error("Failed to fetch size price");
  }

  const buffer = await res.arrayBuffer();

  const message = Proto.SizePrice.decode(new Uint8Array(buffer));
  const data = Proto.SizePrice.toObject(message, { longs: Number });

  if (data.dealer) {
    return {
      client: data.dealer.clientPrice,
      dealer: data.dealer.dealerPrice,
    };
  }

  return {
    client: data.client.clientPrice,
    dealer: 0,
  };
}

/*
================================
PRICE
================================
*/

/* Пересчитывает цену ворот и всего заказа */
function recalculatePrices() {
  for (const gate of orderGates) {
    gate.gatePrice = {
      client: calculateGatePrice(gate, "client"),
      dealer: state.role === "dealer" ? calculateGatePrice(gate, "dealer") : 0,
    };
  }

  state.orderPrice.client = calculateOrderTotal("client");

  state.orderPrice.dealer =
    state.role === "dealer" ? calculateOrderTotal("dealer") : 0;

  renderTotalPrice();
  renderGateList();
}

/* Вычисляет цену одних ворот */
function calculateGatePrice(gate, role) {
  const base = gate.sizePrice?.[role] || 0;

  let total =
    base +
    (base * (gate.liftMarkup?.[role] || 0)) / 100 +
    (base * (gate.cycleMarkup?.[role] || 0)) / 100 +
    (gate.drivePrice?.[role] || 0);

  gate.options.forEach((o) => {
    total += (o.price?.[role] || 0) * (o.amount || 0);
  });

  return Math.round(total);
}

/* Вычисляет общую цену заказа */
function calculateOrderTotal(role) {
  const gates = orderGates.reduce(
    (sum, gate) => sum + (gate.gatePrice?.[role] || 0),
    0,
  );

  const products = state.products.reduce(
    (sum, product) =>
      sum + (product.price?.[role] || 0) * (product.amount || 0),
    0,
  );

  return Math.round(gates + products);
}

/*
================================
RENDER
================================
*/

/* Рендер всего калькулятора */
function renderCalculator() {
  renderCurrentGateForm();
  renderGateList();
  renderTotalPrice();
}

/* Рендер списка опций ворот */
function renderOptionsFromGate() {
  DOM.optionList.innerHTML = "";

  currentGate().options.forEach((item) => {
    const clone = DOM.optionTemplate.content.cloneNode(true);

    const div = clone.querySelector(".option-div");
    const select = div.querySelector("select");
    const input = div.querySelector("input");

    select.value = String(item.id);
    input.value = String(item.amount);

    DOM.optionList.appendChild(clone);
  });
}

/* Рендер формы текущих ворот */
function renderCurrentGateForm() {
  const gate = currentGate();

  DOM.headerGateNumber.textContent = `Ворота №${currentGateIndex + 1}`;

  DOM.width.value = gate.size.width;
  DOM.height.value = gate.size.height;
  DOM.headroom.value = gate.headroom ?? 0;
  DOM.liftType.value = String(gate.liftType ?? "");
  DOM.cycleAmount.value = String(gate.cycleAmount ?? "");
  DOM.drive.value = String(gate.driveId ?? "");
  DOM.colorOut.value = String(gate.colorOutId ?? "");

  renderOptionsFromGate();
}

/* Рендер цен */
function renderTotalPrice() {
  const gate = currentGate();

  if (DOM.gateClientPrice) {
    DOM.gateClientPrice.textContent =
      (gate.gatePrice?.client || 0).toLocaleString() + " руб.";
  }

  if (DOM.totalClientPrice) {
    DOM.totalClientPrice.textContent =
      state.orderPrice.client.toLocaleString() + " руб.";
  }

  if (state.role === "dealer") {
    if (DOM.gateDealerPrice) {
      DOM.gateDealerPrice.textContent =
        (gate.gatePrice?.dealer || 0).toLocaleString() + " руб.";
    }

    if (DOM.totalDealerPrice) {
      DOM.totalDealerPrice.textContent =
        state.orderPrice.dealer.toLocaleString() + " руб.";
    }
  }
}

/* Рендер списка ворот */
function renderGateList() {
  DOM.gateList.innerHTML = "";

  orderGates.forEach((_, index) => {
    const li = document.createElement("li");

    const btn = document.createElement("button");
    btn.type = "button";
    btn.textContent = `Ворота ${index + 1}`;

    if (index === currentGateIndex) {
      btn.classList.add("active");
    }

    btn.addEventListener("click", () => {
      switchGate(index);
    });

    const deleteBtn = document.createElement("button");
    deleteBtn.type = "button";
    deleteBtn.textContent = "✕";
    deleteBtn.classList.add("deleteGate");

    deleteBtn.addEventListener("click", (e) => {
      e.stopPropagation();
      removeGate(index);
    });

    li.appendChild(btn);
    li.appendChild(deleteBtn);

    DOM.gateList.appendChild(li);
  });
}

/*
================================
ORDER
================================
*/

/* Добавляет новые ворота */
async function addGate() {
  const newGate = createGateFromDefaults();

  orderGates.push(newGate);

  currentGateIndex = orderGates.length - 1;

  renderCurrentGateForm();

  await updateCurrentGateSizePrice();

  recalculatePrices();
}

/* Переключает текущие ворота */
function switchGate(index) {
  if (index === currentGateIndex) return;

  parseCurrentGateFromForm();

  currentGateIndex = index;

  renderCurrentGateForm();
  renderTotalPrice();
  renderGateList();
}

/* Удаляет ворота из заказа */
function removeGate(index) {
  if (orderGates.length === 1) return; // нельзя удалить последние ворота

  orderGates.splice(index, 1);

  if (index === currentGateIndex) {
    // если удалили текущие ворота
    currentGateIndex = Math.max(0, index - 1);
  } else if (index < currentGateIndex) {
    // если удалили ворота перед текущими
    currentGateIndex--;
  }

  renderCurrentGateForm();
  recalculatePrices();
}

/*
================================
UTILS
================================
*/

/* Добавляет строку опции */
function addOptionRow() {
  const clone = DOM.optionTemplate.content.cloneNode(true);

  DOM.optionList.appendChild(clone);

  updateForm();
}

/* Добавляет строку товара */
function addProductRow() {
  const clone = DOM.productTemplate.content.cloneNode(true);

  DOM.productList.appendChild(clone);

  updateProducts();
}

/* Удаляет строку опции или товара */
function removeItem(button, selector) {
  const item = button.closest(selector);

  if (item) {
    item.remove();
  }

  if (selector === ".option-div") {
    updateForm();
  }

  if (selector === ".product-div") {
    updateProducts();
  }
}

/* Проверяет количество товара */
function getAmountFromInput(input) {
  const value = Number(input.value);

  if (Number.isNaN(value) || value < 1) {
    return 1;
  }

  return value;
}

/* debounce для ограничений вызовов */
function debounce(fn, delay) {
  let timer;

  return (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delay);
  };
}

/*
================================
ORDER API
================================
*/

/* Преобразует список товаров для API */
function buildProductsList() {
  const map = {};

  state.products.forEach((p) => {
    map[p.id] = p.amount;
  });

  return map;
}

/* Формирует payload ворот для заказа */
function buildGatePayload(gate) {
  const optionsMap = {};

  gate.options.forEach((o) => {
    optionsMap[o.id] = o.amount;
  });

  const gatePrice =
    state.role === "dealer" ? gate.gatePrice.dealer : gate.gatePrice.client;

  return {
    gateType: GateType[gate.gateType],
    width: gate.size.width,
    height: gate.size.height,
    liftTypeId: gate.liftType,
    colorOutId: gate.colorOutId,
    driveId: gate.driveId,
    cycleAmountId: gate.cycleAmount,
    options: optionsMap,
    gatePrice,
    headroom: gate.headroom,
  };
}

/* Отправляет заказ на сервер */
async function placeOrder() {
  const payload = {
    orderGates: orderGates.map(buildGatePayload),
    products: buildProductsList(),
  };

  const err = Proto.OrderRequest.verify(payload);
  if (err) throw new Error(err);

  const message = Proto.OrderRequest.create(payload);
  const buffer = Proto.OrderRequest.encode(message).finish();

  const res = await fetch("/orders", {
    method: "POST",
    headers: { "Content-Type": "application/x-protobuf" },
    body: buffer,
  });

  if (!res.ok) {
    throw new Error("Order request failed");
  }

  window.location.href = "/orders";
}

/*
================================
START
================================
*/

/* Точка входа приложения */
document.addEventListener("DOMContentLoaded", async () => {
  await initProtobuf();
  await initCalculator();
});
