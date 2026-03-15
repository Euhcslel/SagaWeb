"use strict";

/*
==============================
PRICE CALCULATION
==============================
*/

function calcGate(gate, role) {
  const base = gate.sizePrice?.[role] || 0;

  let total =
    base +
    (base * (gate.liftMarkup?.[role] || 0)) / 100 +
    (base * (gate.cycleMarkup?.[role] || 0)) / 100 +
    (gate.drivePrice?.[role] || 0);

  gate.options?.forEach((o) => {
    total += (o.price?.[role] || 0) * (o.amount || 0);
  });

  return total;
}

/*
==============================
SIZE PRICE API
==============================
*/

async function fetchSizePrice(width, height, gateType, Proto) {
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
      client: data.dealer.clientPrice / 100,
      dealer: data.dealer.dealerPrice / 100,
    };
  }

  return {
    client: data.client.clientPrice / 100,
    dealer: 0,
  };
}

/*
==============================
ORDER PAYLOAD
==============================
*/

function buildGatePayload(gate, GateType) {
  const options = {};

  gate.options?.forEach((o) => {
    options[o.id] = o.amount;
  });

  return {
    gateType: GateType[gate.gateType],

    width: gate.size.width,
    height: gate.size.height,

    liftTypeId: gate.liftType,
    colorOutId: gate.colorOutId,
    driveId: gate.driveId,
    cycleAmountId: gate.cycleAmount,

    options,

    headroom: gate.headroom,

    drive: {
      industrial: {
        driveId: 7,
      },
    },
  };
}

/*
==============================
UTILS
==============================
*/

function num(value) {
  const n = Number(value);
  return Number.isNaN(n) ? 0 : n;
}

function clamp(value, min, max) {
  return Math.max(min, Math.min(max, num(value)));
}

function amount(input) {
  const v = num(input?.value);
  return v < 1 ? 1 : v;
}

function dealer(value, role) {
  return role === "dealer" ? num(value) : 0;
}

function markup(item, role) {
  return {
    client: num(item?.RetailMarkup),
    dealer: dealer(item?.WholesaleMarkup, role),
  };
}

function price(item, role) {
  return {
    client: num(item?.RetailPrice),
    dealer: dealer(item?.WholesalePrice, role),
  };
}

function fmt(value) {
  return (
    value.toLocaleString("ru-RU", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }) + " руб."
  );
}

function debounce(fn, delay) {
  let timer;

  return (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delay);
  };
}

/* Не стандарт */
