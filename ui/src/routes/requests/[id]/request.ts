import type { components } from "$lib/oapi/gen/types";

type Request = components['schemas']['Request'];

export type getRequestFunc = () => Request;
