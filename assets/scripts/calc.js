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
    gateType: document.getElementById("gateType"),

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
  initUserRole();
  initGateType();
  initEventListeners();

  applyGateConfig(state.gateType);

  const gate = createGateFromDefaults();
  orderGates.push(gate);
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
function initGateType() {
  state.gateType = DOM.gateType?.value || "res";
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

  DOM.gateType.addEventListener("change", changeGateType);
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
CONFIG
================================
*/

/* Возвращает конфигурацию ворот по типу */
function getCurrentConfig(type = state.gateType) {
  return window.gateConfigs?.[type] || null;
}

/* Возвращает список приводов из конфигурации ворот */
function getDrivesFromConfig(cfg) {
  if (!cfg) return [];
  return cfg.IndustrialDrives?.length
    ? cfg.IndustrialDrives
    : cfg.ResidentialDrives || [];
}

/* Возвращает дилерскую цену, если пользователь дилер */
function getDealerPrice(value) {
  return state.role === "dealer" ? Number(value || 0) : 0;
}

/* Заполняет select элемент списком options */
function fillSelect(selectEl, items, buildOption) {
  selectEl.innerHTML = "";

  items.forEach((item) => {
    const option = buildOption(item);
    selectEl.appendChild(option);
  });
}

/* Создаёт option для select с наценкой (%) */
function buildMarkupOption(item) {
  const option = document.createElement("option");
  option.value = String(item.ID);
  option.textContent = item.Name;
  option.dataset.clientMarkup = String(item.RetailMarkup || 0);
  option.dataset.dealerMarkup = String(getDealerPrice(item.WholesaleMarkup));
  return option;
}

/* Создаёт option для выбора цвета */
function buildColorOption(item) {
  const option = document.createElement("option");
  option.value = String(item.ID);
  option.textContent = item.Code;
  return option;
}

/* Создаёт option для выбора привода */
function buildDriveOption(item) {
  const option = document.createElement("option");
  option.value = String(item.ID);
  option.textContent = item.Name;
  option.dataset.clientPrice = String(item.RetailPrice || 0);
  option.dataset.dealerPrice = String(getDealerPrice(item.WholesalePrice));
  return option;
}

/* Создаёт option с фиксированной ценой */
function buildPriceOption(item) {
  const option = document.createElement("option");
  option.value = String(item.ID);
  option.textContent = item.Name;
  option.dataset.clientPrice = String(item.RetailPrice || 0);
  option.dataset.dealerPrice = String(getDealerPrice(item.WholesalePrice));
  return option;
}

/* Заполняет select внутри template */
function fillTemplateSelect(templateEl, items, buildOption) {
  const select = templateEl.content.querySelector("select");
  if (!select) return;

  select.innerHTML = "";
  items.forEach((item) => {
    select.appendChild(buildOption(item));
  });
}

/* Применяет ограничения размеров ворот из конфигурации */
function applySizeParams(cfg) {
  DOM.width.min = String(cfg.WidthParams?.MinValue ?? 0);
  DOM.width.max = String(cfg.WidthParams?.MaxValue ?? 0);

  DOM.height.min = String(cfg.HeightParams?.MinValue ?? 0);
  DOM.height.max = String(cfg.HeightParams?.MaxValue ?? 0);

  DOM.width.value = String(
    clampNumber(DOM.width.value, DOM.width.min, DOM.width.max),
  );
  DOM.height.value = String(
    clampNumber(DOM.height.value, DOM.height.min, DOM.height.max),
  );
}

/* Применяет конфигурацию ворот к форме */
function applyGateConfig(type) {
  const cfg = getCurrentConfig(type);
  if (!cfg) return;

  fillSelect(DOM.liftType, cfg.LiftTypes || [], buildMarkupOption);
  fillSelect(DOM.cycleAmount, cfg.CycleAmounts || [], buildMarkupOption);
  fillSelect(DOM.colorOut, cfg.Colors || [], buildColorOption);
  fillSelect(DOM.drive, getDrivesFromConfig(cfg), buildDriveOption);

  fillTemplateSelect(DOM.optionTemplate, cfg.Options || [], buildPriceOption);
  fillTemplateSelect(DOM.productTemplate, cfg.Products || [], buildPriceOption);

  applySizeParams(cfg);
}

/* Нормализует параметры ворот согласно конфигурации */
function normalizeGateByConfig(gate) {
  const cfg = getCurrentConfig(gate.gateType);
  if (!cfg) return;

  gate.size.width = clampNumber(
    gate.size.width,
    cfg.WidthParams?.MinValue,
    cfg.WidthParams?.MaxValue,
  );

  gate.size.height = clampNumber(
    gate.size.height,
    cfg.HeightParams?.MinValue,
    cfg.HeightParams?.MaxValue,
  );

  const lift = findById(cfg.LiftTypes, gate.liftType) || cfg.LiftTypes?.[0];
  gate.liftType = Number(lift?.ID || 0);
  gate.liftMarkup = {
    client: Number(lift?.RetailMarkup || 0),
    dealer: getDealerPrice(lift?.WholesaleMarkup),
  };

  const cycle =
    findById(cfg.CycleAmounts, gate.cycleAmount) || cfg.CycleAmounts?.[0];
  gate.cycleAmount = Number(cycle?.ID || 0);
  gate.cycleMarkup = {
    client: Number(cycle?.RetailMarkup || 0),
    dealer: getDealerPrice(cycle?.WholesaleMarkup),
  };

  const drives = getDrivesFromConfig(cfg);
  const drive = findById(drives, gate.driveId) || drives[0];
  gate.driveId = Number(drive?.ID || 0);
  gate.drivePrice = {
    client: Number(drive?.RetailPrice || 0),
    dealer: getDealerPrice(drive?.WholesalePrice),
  };

  const color = findById(cfg.Colors, gate.colorOutId) || cfg.Colors?.[0];
  gate.colorOutId = Number(color?.ID || 0);

  const validOptions = cfg.Options || [];
  gate.options = gate.options
    .map((o) => {
      const item = findById(validOptions, o.id);
      if (!item) return null;

      return {
        id: Number(item.ID),
        amount: getSafeAmount(o.amount),
        price: {
          client: Number(item.RetailPrice || 0),
          dealer: getDealerPrice(item.WholesalePrice),
        },
      };
    })
    .filter(Boolean);
}

/* Обрабатывает смену типа ворот */
function changeGateType() {
  const gate = currentGate();
  const newType = DOM.gateType.value;

  parseCurrentGateFromForm();

  gate.gateType = newType;
  state.gateType = newType;

  applyGateConfig(newType);
  normalizeGateByConfig(gate);
  syncFormWithGate(gate);

  renderCurrentGateForm();
  updateProducts();
  updateSize();
}

/*
================================
DEFAULT GATE
================================
*/

/* Значение ширины по умолчанию */
function getDefaultWidth() {
  return toNumber(DOM.width.value, 0);
}

/* Значение высоты по умолчанию */
function getDefaultHeight() {
  return toNumber(DOM.height.value, 0);
}

/* Значение headroom по умолчанию */
function getDefaultHeadroom() {
  return toNumber(DOM.headroom.value, 0);
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

  gate.gateType = DOM.gateType.value;
  state.gateType = gate.gateType;

  gate.size.width = toNumber(DOM.width.value, 0);
  gate.size.height = toNumber(DOM.height.value, 0);
  gate.headroom = toNumber(DOM.headroom.value, 0);

  const lift = DOM.liftType.selectedOptions[0];
  gate.liftType = Number(lift?.value || 0);
  gate.liftMarkup = {
    client: Number(lift?.dataset.clientMarkup || 0),
    dealer: Number(lift?.dataset.dealerMarkup || 0),
  };

  const cycle = DOM.cycleAmount.selectedOptions[0];
  gate.cycleAmount = Number(cycle?.value || 0);
  gate.cycleMarkup = {
    client: Number(cycle?.dataset.clientMarkup || 0),
    dealer: Number(cycle?.dataset.dealerMarkup || 0),
  };

  const drive = DOM.drive.selectedOptions[0];
  gate.driveId = Number(drive?.value || 0);
  gate.drivePrice = {
    client: Number(drive?.dataset.clientPrice || 0),
    dealer: Number(drive?.dataset.dealerPrice || 0),
  };

  gate.colorOutId = Number(DOM.colorOut.value || 0);

  gate.options = [];

  DOM.optionList.querySelectorAll(".option-div").forEach((div) => {
    const select = div.querySelector("select");
    const input = div.querySelector("input");
    const option = select.options[select.selectedIndex];

    gate.options.push({
      id: Number(option?.value || 0),
      amount: getAmountFromInput(input),
      price: {
        client: Number(option?.dataset.clientPrice || 0),
        dealer: Number(option?.dataset.dealerPrice || 0),
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
      id: Number(option?.value || 0),
      amount: getAmountFromInput(input),
      price: {
        client: Number(option?.dataset.clientPrice || 0),
        dealer: Number(option?.dataset.dealerPrice || 0),
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

  const gate = currentGate();
  if (!Number.isFinite(gate.size.width) || !Number.isFinite(gate.size.height)) {
    return;
  }

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

  gate.sizePrice = await fetchSizePrice(
    gate.size.width,
    gate.size.height,
    gate.gateType,
  );
}

/*
================================
API
================================
*/

/* Запрашивает цену размера ворот у сервера */
async function fetchSizePrice(width, height, gateType) {
  const res = await fetch(
    `/sizes?width=${width}&height=${height}&gateType=${gateType}`,
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

  state.gateType = gate.gateType || state.gateType || "res";
  DOM.gateType.value = state.gateType;

  applyGateConfig(state.gateType);
  normalizeGateByConfig(gate);
  syncFormWithGate(gate);

  DOM.headerGateNumber.textContent = `Ворота №${currentGateIndex + 1}`;
  renderOptionsFromGate();
}

/* Синхронизирует значения формы с объектом ворот */
function syncFormWithGate(gate) {
  DOM.width.value = String(gate.size.width);
  DOM.height.value = String(gate.size.height);
  DOM.headroom.value = String(gate.headroom ?? 0);
  DOM.liftType.value = String(gate.liftType ?? "");
  DOM.cycleAmount.value = String(gate.cycleAmount ?? "");
  DOM.drive.value = String(gate.driveId ?? "");
  DOM.colorOut.value = String(gate.colorOutId ?? "");
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
  parseCurrentGateFromForm();

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
  if (orderGates.length === 1) return;

  orderGates.splice(index, 1);

  if (index === currentGateIndex) {
    currentGateIndex = Math.max(0, index - 1);
  } else if (index < currentGateIndex) {
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

/* Возвращает безопасное количество (минимум 1) */
function getSafeAmount(value) {
  const num = Number(value);
  return Number.isNaN(num) || num < 1 ? 1 : num;
}

/* Преобразует значение в число с fallback */
function toNumber(value, fallback = 0) {
  const num = Number(value);
  return Number.isNaN(num) ? fallback : num;
}

/* Ограничивает число диапазоном min/max */
function clampNumber(value, min, max) {
  let result = toNumber(value, toNumber(min, 0));

  if (min !== undefined && min !== null && min !== "") {
    result = Math.max(result, toNumber(min, result));
  }

  if (max !== undefined && max !== null && max !== "") {
    result = Math.min(result, toNumber(max, result));
  }

  return result;
}

/* Находит элемент в массиве по ID */
function findById(items = [], id) {
  return items.find((item) => Number(item.ID) === Number(id)) || null;
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
  parseCurrentGateFromForm();
  readProductsFromForm();

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
