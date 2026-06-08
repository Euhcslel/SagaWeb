// Entry point for bundling - exports all schemas and runtime helpers
export { create, toBinary, fromBinary } from "@bufbuild/protobuf";

export {
  OrderRequestSchema,
  UpdateProductsRequestSchema,
  GateConfigSchema,
} from "./order_pb.js";

export { SizePriceSchema } from "./prices_pb.js";

export { DocumentsListSchema } from "./documents_pb.js";

export { UpdateOrderStatusRequestSchema } from "./status_pb.js";
