/*
==============================
STATE
==============================
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

const GateType = { ind: 0, res: 1 };

let DOM = {};

/*
==============================
DOM
==============================
*/

function initDomRefs() {
  DOM = {
    isDealer: document.getElementById("isDealer"),

    width: document.getElementById("width"),
    height: document.getElementById("height"),
    headroom: document.getElementById("headroom"),
    liftType: document.getElementById("liftType"),
    cycleAmount: document.getElementById("cycleAmount"),
    drive: document.getElementById("drive"),
    driveType: document.getElementById("driveType"),
    chainLength: document.getElementById("chainLength"),
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

    manualBlock: document.querySelector(".manual"),
    rail: document.getElementById("rail"),
  };
}

/*
==============================
INIT
==============================
*/

async function initProtobuf() {
  const pricesRoot = await protobuf.load("/assets/proto_files/prices.proto");
  Proto.SizePrice = pricesRoot.lookupType("proto.SizePrice");

  const orderRoot = await protobuf.load("/assets/proto_files/order.proto");
  Proto.OrderRequest = orderRoot.lookupType("proto.OrderRequest");
}

async function initCalculator() {
  initDomRefs();
  initUserRole();
  initGateType();

  const gate = createGate(state.gateType);
  orderGates.push(gate);

  applyGateConfig(state.gateType, gate);

  initEvents();

  await updateGateSizePrice(gate);

  recalc();
  render();
}

function initUserRole() {
  state.role = DOM.isDealer?.value === "true" ? "dealer" : "client";
}

function initGateType() {
  state.gateType = DOM.gateType?.value || "res";
}

/*
==============================
CONFIG
==============================
*/

function getConfig(type = state.gateType) {
  return window.gateConfigs?.[type];
}

function getDriveListByGateType(cfg, gateType) {
  if (gateType === "ind") return cfg.IndustrialDrives || [];
  return cfg.ResidentialDrives || [];
}

function getDriveTypeLabel(type) {
  if (type === "manual") return "Ручной";
  if (type === "industrial") return "Промышленный";
  if (type === "residential") return "Бытовой";
  return type;
}

function buildOption(item, type) {
  const option = document.createElement("option");

  if (type === "driveType") {
    option.value = item;
    option.textContent = getDriveTypeLabel(item);
    return option;
  }

  option.value = item.ID;

  if (type === "color") {
    option.textContent = item.Code;
    return option;
  }

  option.textContent = item.Name;

  if (type === "markup") {
    option.dataset.clientMarkup = item.RetailMarkup || 0;
    option.dataset.dealerMarkup = dealer(item.WholesaleMarkup);
  }

  if (type === "price") {
    option.dataset.clientPrice = item.RetailPrice || 0;
    option.dataset.dealerPrice = dealer(item.WholesalePrice);
  }

  return option;
}

function fillSelect(el, list, type) {
  if (!el) return;
  el.innerHTML = "";
  list.forEach((i) => el.appendChild(buildOption(i, type)));
}

function fillTemplate(template, list) {
  const select = template.content.querySelector("select");
  select.innerHTML = "";
  list.forEach((i) => select.appendChild(buildOption(i, "price")));
}

function normalizeDriveType(cfg, driveType) {
  if (!cfg?.DriveTypes?.length) return "manual";
  if (cfg.DriveTypes.includes(driveType)) return driveType;
  return cfg.DriveTypes[0];
}

function setDriveSelectByGate(gate) {
  const cfg = getConfig(gate.gateType);
  const list = getDriveListByGateType(cfg, gate.gateType);

  fillSelect(DOM.drive, list, "price");

  if (!list.length) {
    gate.driveId = 0;
    gate.drivePrice = { client: 0, dealer: 0 };
    return;
  }

  const exists = list.some((d) => d.ID === gate.driveId);

  if (!exists) {
    gate.driveId = list[0].ID;
    gate.drivePrice = price(list[0]);
  }

  DOM.drive.value = String(gate.driveId);
}

function applyGateConfig(type, gate = currentGate()) {
  const cfg = getConfig(type);
  if (!cfg || !gate) return;

  fillSelect(DOM.liftType, cfg.LiftTypes, "markup");
  fillSelect(DOM.cycleAmount, cfg.CycleAmounts, "markup");
  fillSelect(DOM.colorOut, cfg.Colors, "color");
  fillSelect(DOM.driveType, cfg.DriveTypes, "driveType");
  if (type === "res") {
    fillSelect(DOM.rail, cfg.Rails, "price");

    const rail = cfg.Rails[0];

    gate.railId = rail.ID;
    gate.railPrice = price(rail);

    DOM.rail.value = String(rail.ID);
  }

  fillTemplate(DOM.optionTemplate, cfg.Options);
  fillTemplate(DOM.productTemplate, cfg.Products);

  DOM.width.min = cfg.WidthParams.MinValue;
  DOM.width.max = cfg.WidthParams.MaxValue;

  DOM.height.min = cfg.HeightParams.MinValue;
  DOM.height.max = cfg.HeightParams.MaxValue;

  gate.driveType = normalizeDriveType(cfg, gate.driveType);

  setDriveSelectByGate(gate);
  syncDriveUi(gate);
}

/*
==============================
GATE
==============================
*/

function createGate(type) {
  const cfg = getConfig(type);

  const lift = cfg.LiftTypes[0];
  const cycle = cfg.CycleAmounts[0];
  const color = cfg.Colors[0];
  const driveType = cfg.DriveTypes?.[0] || "manual";

  const driveList = getDriveListByGateType(cfg, type);
  const drive = driveList[0] || null;

  return {
    gateType: type,

    size: {
      width: cfg.WidthParams.MinValue,
      height: cfg.HeightParams.MinValue,
    },

    headroom: 0,

    sizePrice: { client: 0, dealer: 0 },

    liftType: lift.ID,
    liftMarkup: markup(lift),

    cycleAmount: cycle.ID,
    cycleMarkup: markup(cycle),

    driveType: driveType,
    chainLength: 0,

    driveId: drive ? drive.ID : 0,
    drivePrice:
      driveType === "manual"
        ? { client: 0, dealer: 0 }
        : drive
          ? price(drive)
          : { client: 0, dealer: 0 },

    railId: 0,
    railPrice: { client: 0, dealer: 0 },

    colorOutId: color.ID,

    options: [],

    gatePrice: { client: 0, dealer: 0 },

    _req: 0,
  };
}

function currentGate() {
  return orderGates[currentGateIndex];
}

function normalizeGate(gate) {
  const cfg = getConfig(gate.gateType);

  gate.size.width = clamp(
    gate.size.width,
    cfg.WidthParams.MinValue,
    cfg.WidthParams.MaxValue,
  );
  gate.size.height = clamp(
    gate.size.height,
    cfg.HeightParams.MinValue,
    cfg.HeightParams.MaxValue,
  );
}

/*
==============================
DRIVE AND RAIL
==============================
*/

function calcManualDrivePrice(chainLength) {
  const value = num(chainLength) * 666;
  return {
    client: value,
    dealer: state.role === "dealer" ? value : 0,
  };
}

function syncDriveUi(gate = currentGate()) {
  if (!gate) return;

  DOM.driveType.value = gate.driveType;
  DOM.chainLength.value = gate.chainLength;

  if (gate.driveType === "manual") {
    DOM.manualBlock.style.display = "block";
    DOM.drive.disabled = true;
  } else {
    DOM.manualBlock.style.display = "none";
    DOM.drive.disabled = false;
  }

  if (!DOM.drive.disabled && gate.driveId) {
    DOM.drive.value = String(gate.driveId);
  }
}

function updateDriveStateFromType(gate = currentGate()) {
  if (!gate) return;

  if (gate.driveType === "manual") {
    gate.drivePrice = calcManualDrivePrice(gate.chainLength);
    gate.driveId = 0;
    return;
  }

  const o = DOM.drive.selectedOptions[0];
  if (!o) {
    gate.driveId = 0;
    gate.drivePrice = { client: 0, dealer: 0 };
    return;
  }

  gate.driveId = num(o.value);
  gate.drivePrice = {
    client: num(o.dataset.clientPrice),
    dealer: num(o.dataset.dealerPrice),
  };
}

function syncRailUi(gate) {
  const show = gate.gateType === "res" && gate.driveType === "residential";

  DOM.rail.parentElement.style.display = show ? "block" : "none";
}

/*
==============================
EVENTS
==============================
*/

function initEvents() {
  const sizeUpdate = debounce(async () => {
    const g = currentGate();

    g.size.width = num(DOM.width.value);
    g.size.height = num(DOM.height.value);

    normalizeGate(g);

    await updateGateSizePrice(g);

    recalc();
    render();
  }, 300);

  DOM.width.addEventListener("input", sizeUpdate);
  DOM.height.addEventListener("input", sizeUpdate);

  DOM.headroom.addEventListener("input", () => {
    currentGate().headroom = num(DOM.headroom.value);
  });

  DOM.liftType.addEventListener("change", () => {
    const o = DOM.liftType.selectedOptions[0];
    const g = currentGate();

    g.liftType = num(o.value);
    g.liftMarkup = {
      client: num(o.dataset.clientMarkup),
      dealer: num(o.dataset.dealerMarkup),
    };

    recalc();
    render();
  });

  DOM.cycleAmount.addEventListener("change", () => {
    const o = DOM.cycleAmount.selectedOptions[0];
    const g = currentGate();

    g.cycleAmount = num(o.value);
    g.cycleMarkup = {
      client: num(o.dataset.clientMarkup),
      dealer: num(o.dataset.dealerMarkup),
    };

    recalc();
    render();
  });

  DOM.driveType.addEventListener("change", () => {
    const g = currentGate();

    g.driveType = DOM.driveType.value;

    updateDriveStateFromType(g);
    syncDriveUi(g);

    if (g.gateType === "res" && g.driveType === "residential") {
      const cfg = getConfig(g.gateType);

      fillSelect(DOM.rail, cfg.Rails, "price");

      const rail = cfg.Rails[0];

      g.railId = rail.ID;
      g.railPrice = price(rail);

      DOM.rail.value = String(rail.ID);
    } else {
      g.railId = 0;
      g.railPrice = { client: 0, dealer: 0 };
    }

    syncRailUi(g);

    recalc();
    renderPrices();
  });

  DOM.chainLength.addEventListener("input", () => {
    const g = currentGate();

    g.chainLength = num(DOM.chainLength.value);

    if (g.driveType === "manual") {
      g.drivePrice = calcManualDrivePrice(g.chainLength);
      recalc();
      renderPrices();
    }
  });

  DOM.drive.addEventListener("change", () => {
    const g = currentGate();

    if (g.driveType === "manual") return;

    const o = DOM.drive.selectedOptions[0];

    g.driveId = num(o.value);
    g.drivePrice = {
      client: num(o.dataset.clientPrice),
      dealer: num(o.dataset.dealerPrice),
    };

    recalc();
    render();
  });

  DOM.rail.addEventListener("change", () => {
    const g = currentGate();
    const o = DOM.rail.selectedOptions[0];

    g.railId = num(o.value);

    g.railPrice = {
      client: num(o.dataset.clientPrice),
      dealer: num(o.dataset.dealerPrice),
    };

    recalc();
    renderPrices();
  });

  DOM.colorOut.addEventListener("change", () => {
    currentGate().colorOutId = num(DOM.colorOut.value);
  });

  DOM.optionList.addEventListener("input", updateOptions);
  DOM.productList.addEventListener("input", updateProducts);

  DOM.addGateBtn.onclick = addGate;
  DOM.addOptionBtn.onclick = addOptionRow;
  DOM.addProductBtn.onclick = addProductRow;

  DOM.gateType.onchange = changeGateType;
}

/*
==============================
STATE UPDATE
==============================
*/

function updateOptions() {
  const g = currentGate();

  g.options = [];

  DOM.optionList.querySelectorAll(".option-div").forEach((div) => {
    const s = div.querySelector("select");
    const i = div.querySelector("input");

    const o = s.options[s.selectedIndex];

    g.options.push({
      id: num(o.value),
      amount: amount(i),
      price: {
        client: num(o.dataset.clientPrice),
        dealer: num(o.dataset.dealerPrice),
      },
    });
  });

  recalc();
  render();
}

function updateProducts() {
  state.products = [];

  DOM.productList.querySelectorAll(".product-div").forEach((div) => {
    const s = div.querySelector("select");
    const i = div.querySelector("input");

    const o = s.options[s.selectedIndex];

    state.products.push({
      id: num(o.value),
      amount: amount(i),
      price: {
        client: num(o.dataset.clientPrice),
        dealer: num(o.dataset.dealerPrice),
      },
    });
  });

  recalc();
  render();
}

async function changeGateType() {
  const g = currentGate();
  const newType = DOM.gateType.value;
  const cfg = getConfig(newType);

  g.gateType = newType;
  state.gateType = newType;

  normalizeGate(g);

  const lift = cfg.LiftTypes[0];
  const cycle = cfg.CycleAmounts[0];
  const color = cfg.Colors[0];
  const driveType = cfg.DriveTypes?.[0] || "manual";
  const driveList = getDriveListByGateType(cfg, newType);
  const drive = driveList[0] || null;

  g.liftType = lift.ID;
  g.liftMarkup = markup(lift);

  g.cycleAmount = cycle.ID;
  g.cycleMarkup = markup(cycle);

  g.colorOutId = color.ID;

  g.driveType = driveType;
  g.chainLength = 0;
  g.driveId = driveType === "manual" ? 0 : drive ? drive.ID : 0;
  g.drivePrice =
    driveType === "manual"
      ? calcManualDrivePrice(0)
      : drive
        ? price(drive)
        : { client: 0, dealer: 0 };

  g.railId = 0;
  g.railPrice = { client: 0, dealer: 0 };

  applyGateConfig(g.gateType, g);

  await updateGateSizePrice(g);

  recalc();
  render();
}

/*
==============================
API
==============================
*/

async function updateGateSizePrice(gate) {
  gate._req++;

  const id = gate._req;

  const price = await fetchSizePrice(
    gate.size.width,
    gate.size.height,
    gate.gateType,
  );

  if (id !== gate._req) return;

  gate.sizePrice = price;
}

async function fetchSizePrice(w, h, t) {
  const res = await fetch(`/sizes?width=${w}&height=${h}&gateType=${t}`);

  const buf = await res.arrayBuffer();

  const msg = Proto.SizePrice.decode(new Uint8Array(buf));

  const d = Proto.SizePrice.toObject(msg, { longs: Number });

  if (d.dealer) {
    return {
      client: d.dealer.clientPrice / 100,
      dealer: d.dealer.dealerPrice / 100,
    };
  }

  return { client: d.client.clientPrice / 100, dealer: 0 };
}

/*
==============================
PRICE
==============================
*/

function recalc() {
  for (const g of orderGates) {
    g.gatePrice = {
      client: calcGate(g, "client"),
      dealer: state.role === "dealer" ? calcGate(g, "dealer") : 0,
    };
  }

  state.orderPrice.client = calcOrder("client");
  state.orderPrice.dealer = calcOrder("dealer");
}

function calcGate(g, role) {
  const base = g.sizePrice?.[role] || 0;

  let total =
    base +
    (base * (g.liftMarkup?.[role] || 0)) / 100 +
    (base * (g.cycleMarkup?.[role] || 0)) / 100 +
    (g.drivePrice?.[role] || 0) +
    (g.railPrice?.[role] || 0);

  g.options.forEach((o) => {
    total += (o.price?.[role] || 0) * o.amount;
  });

  return total;
}

function calcOrder(role) {
  const gates = orderGates.reduce((s, g) => s + (g.gatePrice?.[role] || 0), 0);

  const products = state.products.reduce(
    (s, p) => s + (p.price?.[role] || 0) * p.amount,
    0,
  );

  return gates + products;
}

/*
==============================
RENDER
==============================
*/

function render() {
  renderGateForm();
  renderOptions();
  renderProducts();
  renderGateList();
  renderPrices();
}

function renderGateForm() {
  const g = currentGate();
  const cfg = getConfig(g.gateType);

  DOM.headerGateNumber.textContent = `Ворота №${currentGateIndex + 1}`;

  DOM.width.value = g.size.width;
  DOM.height.value = g.size.height;
  DOM.headroom.value = g.headroom;

  fillSelect(DOM.liftType, cfg.LiftTypes, "markup");
  fillSelect(DOM.cycleAmount, cfg.CycleAmounts, "markup");
  fillSelect(DOM.colorOut, cfg.Colors, "color");
  fillSelect(DOM.driveType, cfg.DriveTypes, "driveType");
  setDriveSelectByGate(g);

  if (g.gateType === "res") {
    fillSelect(DOM.rail, cfg.Rails, "price");
  }

  DOM.liftType.value = String(g.liftType);
  DOM.cycleAmount.value = String(g.cycleAmount);
  DOM.colorOut.value = String(g.colorOutId);
  DOM.gateType.value = g.gateType;

  syncDriveUi(g);
  if (g.railId) {
    DOM.rail.value = String(g.railId);
  }
  syncRailUi(g);
}

function renderOptions() {
  const g = currentGate();

  DOM.optionList.innerHTML = "";

  g.options.forEach((o) => {
    const c = DOM.optionTemplate.content.cloneNode(true);

    const s = c.querySelector("select");
    const i = c.querySelector("input");

    s.value = o.id;
    i.value = o.amount;

    DOM.optionList.appendChild(c);
  });
}

function renderProducts() {
  DOM.productList.innerHTML = "";

  state.products.forEach((p) => {
    const c = DOM.productTemplate.content.cloneNode(true);

    const s = c.querySelector("select");
    const i = c.querySelector("input");

    s.value = p.id;
    i.value = p.amount;

    DOM.productList.appendChild(c);
  });
}

function renderGateList() {
  DOM.gateList.innerHTML = "";

  orderGates.forEach((g, i) => {
    const li = document.createElement("li");

    const b = document.createElement("button");
    b.textContent = `Ворота ${i + 1}`;
    b.onclick = () => switchGate(i);

    if (i === currentGateIndex) b.classList.add("active");

    const d = document.createElement("button");
    d.textContent = "✕";
    d.onclick = (e) => {
      e.stopPropagation();
      removeGate(i);
    };

    li.append(b, d);
    DOM.gateList.appendChild(li);
  });
}

function renderPrices() {
  const g = currentGate();

  DOM.gateClientPrice.textContent = fmt(g.gatePrice.client);
  DOM.totalClientPrice.textContent = fmt(state.orderPrice.client);

  if (state.role === "dealer") {
    DOM.gateDealerPrice.textContent = fmt(g.gatePrice.dealer);
    DOM.totalDealerPrice.textContent = fmt(state.orderPrice.dealer);
  }
}

/*
==============================
ORDER ACTIONS
==============================
*/

async function addGate() {
  const g = createGate(state.gateType);

  orderGates.push(g);
  currentGateIndex = orderGates.length - 1;

  applyGateConfig(g.gateType, g);

  await updateGateSizePrice(g);

  recalc();
  render();
}

function switchGate(i) {
  currentGateIndex = i;
  render();
}

function removeGate(i) {
  if (orderGates.length === 1) return;

  orderGates.splice(i, 1);

  if (currentGateIndex >= orderGates.length) {
    currentGateIndex = orderGates.length - 1;
  } else if (currentGateIndex > i) {
    currentGateIndex -= 1;
  }

  recalc();
  render();
}

/*
==============================
UTIL
==============================
*/

function addOptionRow() {
  const cfg = getConfig(currentGate().gateType);
  const o = cfg.Options[0];

  currentGate().options.push({
    id: o.ID,
    amount: 1,
    price: price(o),
  });

  recalc();
  render();
}

function addProductRow() {
  const cfg = getConfig(currentGate().gateType);
  const p = cfg.Products[0];

  state.products.push({
    id: p.ID,
    amount: 1,
    price: price(p),
  });

  recalc();
  render();
}

function removeItem(btn, sel) {
  const el = btn.closest(sel);
  const list = [...el.parentElement.children];
  const i = list.indexOf(el);

  if (sel === ".option-div") currentGate().options.splice(i, 1);
  if (sel === ".product-div") state.products.splice(i, 1);

  recalc();
  render();
}

function markup(i) {
  return { client: num(i.RetailMarkup), dealer: dealer(i.WholesaleMarkup) };
}

function price(i) {
  return { client: num(i.RetailPrice), dealer: dealer(i.WholesalePrice) };
}

function dealer(v) {
  return state.role === "dealer" ? num(v) : 0;
}

function num(v) {
  const n = Number(v);
  return Number.isNaN(n) ? 0 : n;
}

function clamp(v, min, max) {
  return Math.max(min, Math.min(max, num(v)));
}

function amount(input) {
  const v = num(input?.value);
  return v < 1 ? 1 : v;
}

function fmt(v) {
  return v.toLocaleString("ru-RU", { minimumFractionDigits: 2 }) + " руб.";
}

function debounce(fn, ms) {
  let t;
  return (...a) => {
    clearTimeout(t);
    t = setTimeout(() => fn(...a), ms);
  };
}

/*
==============================
ORDER API
==============================
*/

function buildProductsList() {
  const m = {};
  state.products.forEach((p) => {
    m[p.id] = p.amount;
  });
  return m;
}

function buildDrivePayload(g) {
  if (g.driveType === "manual") {
    return {
      manual: {
        chainLength: g.chainLength,
      },
    };
  }

  if (g.driveType === "industrial") {
    return {
      industrial: {
        driveId: g.driveId,
      },
    };
  }

  if (g.driveType === "residential") {
    return {
      residential: {
        driveId: g.driveId,
        railId: g.railId,
      },
    };
  }

  return null;
}

function buildGatePayload(g) {
  const opts = {};
  g.options.forEach((o) => {
    opts[o.id] = o.amount;
  });

  return {
    gateType: GateType[g.gateType],
    width: g.size.width,
    height: g.size.height,
    liftTypeId: g.liftType,
    colorOutId: g.colorOutId,
    cycleAmountId: g.cycleAmount,
    options: opts,
    headroom: g.headroom,
    drive: buildDrivePayload(g),
    price: Math.round(g.gatePrice.dealer * 100),
  };
}

async function placeOrder() {
  const payload = {
    orderGates: orderGates.map(buildGatePayload),
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

/*
==============================
START
==============================
*/

document.addEventListener("DOMContentLoaded", async () => {
  await initProtobuf();
  await initCalculator();
});
