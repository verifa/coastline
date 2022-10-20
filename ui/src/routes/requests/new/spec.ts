import type { SchemaObject } from "openapi-typescript";

export interface Property {
    name: string;
    schema: SchemaObject;
    is_required: boolean;
};

export function propFromSchema(name: string, parentObj: SchemaObject, obj: SchemaObject): Property {
    return {
        name: name,
        schema: obj,
        is_required: (parentObj.required || []).includes(name)
    }

}

export function getInitialPropValue(prop: SchemaObject) {
    // Basic types:
    // Ref: https://swagger.io/docs/specification/data-models/data-types/
    //   string (this includes dates and files)
    //   number
    //   integer
    //   boolean
    //   array
    //   object
    switch (prop.type) {
        case 'string': {
            return prop.default || '';
        }
        case 'number': {
            return prop.default || 0;
        }
        case 'integer': {
            return prop.default || 0;
        }
        case 'boolean': {
            return prop.default || false;
        }
        case 'array': {
            return prop.default || [];
        }
        case 'object': {
            return prop.default || {};
        }
        default: {
            console.log('error: unsupported spec type: ', prop.type);
            return null
        }
    }
}