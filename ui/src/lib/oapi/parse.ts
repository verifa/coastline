import type { OpenAPI3, ReferenceObject, SchemaObject } from 'openapi-typescript'

export interface RequestSpec {
    type: string
    spec: object
}
// import { ReferenceObject } from 'openapi-typescript'

export function getRequestSpecs(spec: OpenAPI3): RequestSpec[] {
    const schemas = spec.components?.schemas
    if (!schemas) {
        return []
    }
    let requestSpecs: RequestSpec[] = []
    for (const name in schemas) {
        const def = schemas[name]

        if ("$ref" in def) {
            // TODO: handle $ref
        } else {
            const obj: SchemaObject = def;

            if (obj && obj.properties && "spec" in obj.properties) {
                const specObj = obj.properties["spec"]
                if ("$ref" in specObj) {
                    // TODO: handle $ref
                } else {
                    requestSpecs.push({
                        type: name,
                        spec: specObj
                    })
                }
            }

        }

    }
    return requestSpecs
}

// function isRef(obj: ReferenceObject | SchemaObject): boolean {
//     return "$ref" in obj
// }