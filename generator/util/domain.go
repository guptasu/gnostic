// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/googleapis/openapi-compiler/jsonschema"
)

// models a collection of types that is defined by a schema
type Domain struct {
	TypeModels         map[string]*TypeModel   // models of the types in the domain
	Prefix             string                  // type prefix to use
	Schema             *jsonschema.Schema      // top-level schema
	PatternNames       map[string]string       // a configured mapping from patterns to property names
	ObjectTypeRequests map[string]*TypeRequest // anonymous types implied by type instantiation
	MapTypeRequests    map[string]string       // "NamedObject" types that will be used to implement ordered maps
}

func NewDomain(schema *jsonschema.Schema) *Domain {
	cc := &Domain{}
	cc.TypeModels = make(map[string]*TypeModel, 0)
	cc.PatternNames = make(map[string]string, 0)
	cc.ObjectTypeRequests = make(map[string]*TypeRequest, 0)
	cc.MapTypeRequests = make(map[string]string, 0)
	cc.Schema = schema
	return cc
}

// Returns a capitalized name to use for a generated type
func (domain *Domain) typeNameForStub(stub string) string {
	return domain.Prefix + strings.ToUpper(stub[0:1]) + stub[1:len(stub)]
}

// Returns a capitalized name to use for a generated type based on a JSON reference
func (domain *Domain) typeNameForReference(reference string) string {
	parts := strings.Split(reference, "/")
	first := parts[0]
	last := parts[len(parts)-1]
	if first == "#" {
		return domain.typeNameForStub(last)
	} else {
		return "Schema"
	}
}

// Returns a property name to use for a JSON reference
func (domain *Domain) propertyNameForReference(reference string) *string {
	parts := strings.Split(reference, "/")
	first := parts[0]
	last := parts[len(parts)-1]
	if first == "#" {
		return &last
	} else {
		return nil
	}
	return nil
}

// Determines the item type for arrays defined by a schema
func (domain *Domain) arrayItemTypeForSchema(propertyName string, schema *jsonschema.Schema) string {
	// default
	itemTypeName := "Any"

	if schema.Items != nil {

		if schema.Items.SchemaArray != nil {

			if len(*(schema.Items.SchemaArray)) > 0 {
				ref := (*schema.Items.SchemaArray)[0].Ref
				if ref != nil {
					itemTypeName = domain.typeNameForReference(*ref)
				} else {
					types := (*schema.Items.SchemaArray)[0].Type
					if types == nil {
						// do nothing
					} else if (types.StringArray != nil) && len(*(types.StringArray)) == 1 {
						itemTypeName = (*types.StringArray)[0]
					} else if (types.StringArray != nil) && len(*(types.StringArray)) > 1 {
						itemTypeName = fmt.Sprintf("%+v", types.StringArray)
					} else if types.String != nil {
						itemTypeName = *(types.String)
					} else {
						itemTypeName = "UNKNOWN"
					}
				}
			}

		} else if schema.Items.Schema != nil {
			types := schema.Items.Schema.Type

			if schema.Items.Schema.Ref != nil {
				itemTypeName = domain.typeNameForReference(*schema.Items.Schema.Ref)
			} else if schema.Items.Schema.OneOf != nil {
				// this type is implied by the "oneOf"
				itemTypeName = domain.typeNameForStub(propertyName + "Item")
				domain.ObjectTypeRequests[itemTypeName] =
					NewTypeRequest(itemTypeName, propertyName, schema.Items.Schema)
			} else if types == nil {
				// do nothing
			} else if (types.StringArray != nil) && len(*(types.StringArray)) == 1 {
				itemTypeName = (*types.StringArray)[0]
			} else if (types.StringArray != nil) && len(*(types.StringArray)) > 1 {
				itemTypeName = fmt.Sprintf("%+v", types.StringArray)
			} else if types.String != nil {
				itemTypeName = *(types.String)
			} else {
				itemTypeName = "UNKNOWN"
			}
		}

	}
	return itemTypeName
}

func (domain *Domain) buildTypeProperties(typeModel *TypeModel, schema *jsonschema.Schema) {
	if schema.Properties != nil {
		for _, pair := range *(schema.Properties) {
			propertyName := pair.Name
			propertySchema := pair.Value
			if propertySchema.Ref != nil {
				// the property schema is a reference, so we will add a property with the type of the referenced schema
				propertyTypeName := domain.typeNameForReference(*(propertySchema.Ref))
				typeProperty := NewTypeProperty()
				typeProperty.Name = propertyName
				typeProperty.Type = propertyTypeName
				typeModel.addProperty(typeProperty)
			} else if propertySchema.Type != nil {
				// the property schema specifies a type, so add a property with the specified type
				if propertySchema.TypeIs("string") {
					typeProperty := NewTypePropertyWithNameAndType(propertyName, "string")
					if propertySchema.Description != nil {
						typeProperty.Description = *propertySchema.Description
					}
					typeModel.addProperty(typeProperty)
				} else if propertySchema.TypeIs("boolean") {
					typeProperty := NewTypePropertyWithNameAndType(propertyName, "bool")
					if propertySchema.Description != nil {
						typeProperty.Description = *propertySchema.Description
					}
					typeModel.addProperty(typeProperty)
				} else if propertySchema.TypeIs("number") {
					typeProperty := NewTypePropertyWithNameAndType(propertyName, "float")
					if propertySchema.Description != nil {
						typeProperty.Description = *propertySchema.Description
					}
					typeModel.addProperty(typeProperty)
				} else if propertySchema.TypeIs("integer") {
					typeProperty := NewTypePropertyWithNameAndType(propertyName, "int")
					if propertySchema.Description != nil {
						typeProperty.Description = *propertySchema.Description
					}
					typeModel.addProperty(typeProperty)
				} else if propertySchema.TypeIs("object") {
					// the property has an "anonymous" object schema, so define a new type for it and request its creation
					anonymousObjectTypeName := domain.typeNameForStub(propertyName)
					domain.ObjectTypeRequests[anonymousObjectTypeName] =
						NewTypeRequest(anonymousObjectTypeName, propertyName, propertySchema)
					// add a property with the type of the requested type
					typeProperty := NewTypePropertyWithNameAndType(propertyName, anonymousObjectTypeName)
					if propertySchema.Description != nil {
						typeProperty.Description = *propertySchema.Description
					}
					typeModel.addProperty(typeProperty)
				} else if propertySchema.TypeIs("array") {
					// the property has an array type, so define it as a a repeated property of the specified type
					propertyTypeName := domain.arrayItemTypeForSchema(propertyName, propertySchema)
					typeProperty := NewTypePropertyWithNameAndType(propertyName, propertyTypeName)
					typeProperty.Repeated = true
					if propertySchema.Description != nil {
						typeProperty.Description = *propertySchema.Description
					}
					typeModel.addProperty(typeProperty)
				} else {
					log.Printf("ignoring %+v, which has an unsupported property type '%+v'", propertyName, propertySchema.Type)
				}
			} else if propertySchema.IsEmpty() {
				// an empty schema can contain anything, so add an accessor for a generic object
				typeName := "Any"
				typeProperty := NewTypePropertyWithNameAndType(propertyName, typeName)
				typeModel.addProperty(typeProperty)
			} else if propertySchema.OneOf != nil {
				anonymousObjectTypeName := domain.typeNameForStub(propertyName + "Item")
				domain.ObjectTypeRequests[anonymousObjectTypeName] =
					NewTypeRequest(anonymousObjectTypeName, propertyName, propertySchema)
				typeProperty := NewTypePropertyWithNameAndType(propertyName, anonymousObjectTypeName)
				typeModel.addProperty(typeProperty)
			} else if propertySchema.AnyOf != nil {
				anonymousObjectTypeName := domain.typeNameForStub(propertyName + "Item")
				domain.ObjectTypeRequests[anonymousObjectTypeName] =
					NewTypeRequest(anonymousObjectTypeName, propertyName, propertySchema)
				typeProperty := NewTypePropertyWithNameAndType(propertyName, anonymousObjectTypeName)
				typeModel.addProperty(typeProperty)
			} else {
				log.Printf("ignoring %s.%s, which has an unrecognized schema:\n%+v", typeModel.Name, propertyName, propertySchema.String())
			}
		}
	}
}

func (domain *Domain) buildTypeRequirements(typeModel *TypeModel, schema *jsonschema.Schema) {
	if schema.Required != nil {
		typeModel.Required = (*schema.Required)
	}
}

func (domain *Domain) buildPatternPropertyAccessors(typeModel *TypeModel, schema *jsonschema.Schema) {
	if schema.PatternProperties != nil {
		typeModel.OpenPatterns = make([]string, 0)
		for _, pair := range *(schema.PatternProperties) {
			propertyPattern := pair.Name
			propertySchema := pair.Value
			typeModel.OpenPatterns = append(typeModel.OpenPatterns, propertyPattern)
			typeName := "Any"
			propertyName := domain.PatternNames[propertyPattern]
			if propertySchema.Ref != nil {
				typeName = domain.typeNameForReference(*propertySchema.Ref)
			}
			propertyTypeName := fmt.Sprintf("Named%s", typeName)
			property := NewTypePropertyWithNameTypeAndPattern(propertyName, propertyTypeName, propertyPattern)
			property.Implicit = true
			property.MapType = typeName
			property.Repeated = true
			domain.MapTypeRequests[property.MapType] = property.MapType
			typeModel.addProperty(property)
		}
	}
}

func (domain *Domain) buildAdditionalPropertyAccessors(typeModel *TypeModel, schema *jsonschema.Schema) {
	if schema.AdditionalProperties != nil {
		if schema.AdditionalProperties.Boolean != nil {
			if *schema.AdditionalProperties.Boolean == true {
				typeModel.Open = true
				propertyName := "additionalProperties"
				typeName := "NamedAny"
				property := NewTypePropertyWithNameAndType(propertyName, typeName)
				property.Implicit = true
				property.MapType = "Any"
				property.Repeated = true
				domain.MapTypeRequests[property.MapType] = property.MapType
				typeModel.addProperty(property)
				return
			}
		} else if schema.AdditionalProperties.Schema != nil {
			typeModel.Open = true
			schema := schema.AdditionalProperties.Schema
			if schema.Ref != nil {
				propertyName := "additionalProperties"
				mapType := domain.typeNameForReference(*schema.Ref)
				typeName := fmt.Sprintf("Named%s", mapType)
				property := NewTypePropertyWithNameAndType(propertyName, typeName)
				property.Implicit = true
				property.MapType = mapType
				property.Repeated = true
				domain.MapTypeRequests[property.MapType] = property.MapType
				typeModel.addProperty(property)
				return
			} else if schema.Type != nil {
				typeName := *schema.Type.String
				if typeName == "string" {
					propertyName := "additionalProperties"
					typeName := "NamedString"
					property := NewTypePropertyWithNameAndType(propertyName, typeName)
					property.Implicit = true
					property.MapType = "string"
					property.Repeated = true
					domain.MapTypeRequests[property.MapType] = property.MapType
					typeModel.addProperty(property)
					return
				} else if typeName == "array" {
					if schema.Items != nil {
						itemType := *schema.Items.Schema.Type.String
						if itemType == "string" {
							propertyName := "additionalProperties"
							typeName := "NamedStringArray"
							property := NewTypePropertyWithNameAndType(propertyName, typeName)
							property.Implicit = true
							property.MapType = "StringArray"
							property.Repeated = true
							domain.MapTypeRequests[property.MapType] = property.MapType
							typeModel.addProperty(property)
							return
						}
					}
				}
			} else if schema.OneOf != nil {
				propertyTypeName := domain.typeNameForStub(typeModel.Name + "Item")
				propertyName := "additionalProperties"
				typeName := fmt.Sprintf("Named%s", propertyTypeName)
				property := NewTypePropertyWithNameAndType(propertyName, typeName)
				property.Implicit = true
				property.MapType = propertyTypeName
				property.Repeated = true
				domain.MapTypeRequests[property.MapType] = property.MapType
				typeModel.addProperty(property)

				domain.ObjectTypeRequests[propertyTypeName] =
					NewTypeRequest(propertyTypeName, propertyName, schema)
			}
		}
	}
}

func (domain *Domain) buildOneOfAccessors(typeModel *TypeModel, schema *jsonschema.Schema) {
	oneOfs := schema.OneOf
	if oneOfs == nil {
		return
	}
	typeModel.Open = true
	typeModel.OneOfWrapper = true
	for _, oneOf := range *oneOfs {
		if oneOf.Ref != nil {
			ref := *oneOf.Ref
			typeName := domain.typeNameForReference(ref)
			propertyName := domain.propertyNameForReference(ref)

			if propertyName != nil {
				typeProperty := NewTypePropertyWithNameAndType(*propertyName, typeName)
				typeModel.addProperty(typeProperty)
			}
		} else if oneOf.Type != nil && oneOf.Type.String != nil && *oneOf.Type.String == "boolean" {
			typeProperty := NewTypePropertyWithNameAndType("boolean", "bool")
			typeModel.addProperty(typeProperty)
		} else {
			log.Printf("Unsupported oneOf:\n%+v", oneOf.String())
		}

	}
}

func schemaIsContainedInArray(s1 *jsonschema.Schema, s2 *jsonschema.Schema) bool {
	if s2.TypeIs("array") {
		if s2.Items.Schema != nil {
			if s1.IsEqual(s2.Items.Schema) {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
}

func (domain *Domain) addAnonymousAccessorForSchema(
	typeModel *TypeModel,
	schema *jsonschema.Schema,
	repeated bool) {
	ref := schema.Ref
	if ref != nil {
		typeName := domain.typeNameForReference(*ref)
		propertyName := domain.propertyNameForReference(*ref)
		if propertyName != nil {
			property := NewTypePropertyWithNameAndType(*propertyName, typeName)
			property.Repeated = true
			typeModel.addProperty(property)
			typeModel.IsItemArray = true
		}
	} else {
		typeName := "string"
		propertyName := "value"
		property := NewTypePropertyWithNameAndType(propertyName, typeName)
		property.Repeated = true
		typeModel.addProperty(property)
		typeModel.IsStringArray = true
	}
}

func (domain *Domain) buildAnyOfAccessors(typeModel *TypeModel, schema *jsonschema.Schema) {
	anyOfs := schema.AnyOf
	if anyOfs == nil {
		return
	}
	if len(*anyOfs) == 2 {
		if schemaIsContainedInArray((*anyOfs)[0], (*anyOfs)[1]) {
			log.Printf("ARRAY OF %+v", (*anyOfs)[0].String())
			schema := (*anyOfs)[0]
			domain.addAnonymousAccessorForSchema(typeModel, schema, true)
		} else if schemaIsContainedInArray((*anyOfs)[1], (*anyOfs)[0]) {
			log.Printf("ARRAY OF %+v", (*anyOfs)[1].String())
			schema := (*anyOfs)[1]
			domain.addAnonymousAccessorForSchema(typeModel, schema, true)
		} else {
			for _, anyOf := range *anyOfs {
				ref := anyOf.Ref
				if ref != nil {
					typeName := domain.typeNameForReference(*ref)
					propertyName := domain.propertyNameForReference(*ref)
					if propertyName != nil {
						property := NewTypePropertyWithNameAndType(*propertyName, typeName)
						typeModel.addProperty(property)
					}
				} else {
					typeName := "bool"
					propertyName := "boolean"
					property := NewTypePropertyWithNameAndType(propertyName, typeName)
					typeModel.addProperty(property)
				}
			}
		}
	} else {
		log.Printf("Unhandled anyOfs:\n%s", schema.String())
	}
}

func (domain *Domain) buildDefaultAccessors(typeModel *TypeModel, schema *jsonschema.Schema) {
	typeModel.Open = true
	propertyName := "additionalProperties"
	typeName := "NamedAny"
	property := NewTypePropertyWithNameAndType(propertyName, typeName)
	property.MapType = "Any"
	property.Repeated = true
	domain.MapTypeRequests[property.MapType] = property.MapType
	typeModel.addProperty(property)
}

func (domain *Domain) buildTypeForDefinition(
	typeName string,
	propertyName string,
	schema *jsonschema.Schema) *TypeModel {
	if (schema.Type == nil) || (*schema.Type.String == "object") {
		return domain.buildTypeForDefinitionObject(typeName, propertyName, schema)
	} else {
		return nil
	}
}

func (domain *Domain) buildTypeForDefinitionObject(
	typeName string,
	propertyName string,
	schema *jsonschema.Schema) *TypeModel {
	typeModel := NewTypeModel()
	typeModel.Name = typeName
	if schema.IsEmpty() {
		domain.buildDefaultAccessors(typeModel, schema)
	} else {
		if schema.Description != nil {
			typeModel.Description = *schema.Description
		}
		domain.buildTypeProperties(typeModel, schema)
		domain.buildTypeRequirements(typeModel, schema)
		domain.buildPatternPropertyAccessors(typeModel, schema)
		domain.buildAdditionalPropertyAccessors(typeModel, schema)
		domain.buildOneOfAccessors(typeModel, schema)
		domain.buildAnyOfAccessors(typeModel, schema)
	}
	return typeModel
}

func (domain *Domain) Build() {
	// create a type for the top-level schema
	typeName := domain.Prefix + "Document"
	typeModel := NewTypeModel()
	typeModel.Name = typeName
	domain.buildTypeProperties(typeModel, domain.Schema)
	domain.buildTypeRequirements(typeModel, domain.Schema)
	domain.buildPatternPropertyAccessors(typeModel, domain.Schema)
	domain.buildAdditionalPropertyAccessors(typeModel, domain.Schema)
	domain.buildOneOfAccessors(typeModel, domain.Schema)
	domain.buildAnyOfAccessors(typeModel, domain.Schema)
	domain.TypeModels[typeName] = typeModel

	// create a type for each object defined in the schema
	if domain.Schema.Definitions != nil {
		for _, pair := range *(domain.Schema.Definitions) {
			definitionName := pair.Name
			definitionSchema := pair.Value
			typeName := domain.typeNameForStub(definitionName)
			typeModel := domain.buildTypeForDefinition(typeName, definitionName, definitionSchema)
			if typeModel != nil {
				domain.TypeModels[typeName] = typeModel
			}
		}
	}
	// iterate over anonymous object types to be instantiated and generate a type for each
	for typeName, typeRequest := range domain.ObjectTypeRequests {
		domain.TypeModels[typeRequest.Name] =
			domain.buildTypeForDefinitionObject(typeName, typeRequest.PropertyName, typeRequest.Schema)
	}

	// iterate over map item types to be instantiated and generate a type for each
	mapTypeNames := make([]string, 0)
	for mapTypeName, _ := range domain.MapTypeRequests {
		mapTypeNames = append(mapTypeNames, mapTypeName)
	}
	sort.Strings(mapTypeNames)

	for _, mapTypeName := range mapTypeNames {
		typeName := "Named" + strings.Title(mapTypeName)
		typeModel := NewTypeModel()
		typeModel.Name = typeName
		typeModel.Description = fmt.Sprintf(
			"Automatically-generated message used to represent maps of %s as ordered (name,value) pairs.",
			mapTypeName)
		typeModel.IsPair = true
		typeModel.PairValueType = mapTypeName

		nameProperty := NewTypeProperty()
		nameProperty.Name = "name"
		nameProperty.Type = "string"
		nameProperty.Description = "Map key"
		typeModel.addProperty(nameProperty)

		valueProperty := NewTypeProperty()
		valueProperty.Name = "value"
		valueProperty.Type = mapTypeName
		valueProperty.Description = "Mapped value"
		typeModel.addProperty(valueProperty)

		domain.TypeModels[typeName] = typeModel
	}

	// add a type for string arrays
	stringArrayType := NewTypeModel()
	stringArrayType.Name = "StringArray"
	stringProperty := NewTypeProperty()
	stringProperty.Name = "value"
	stringProperty.Type = "string"
	stringProperty.Repeated = true
	stringArrayType.addProperty(stringProperty)
	domain.TypeModels[stringArrayType.Name] = stringArrayType

	// add a type for "Any"
	anyType := NewTypeModel()
	anyType.Name = "Any"
	anyType.Open = true
	anyType.IsBlob = true
	valueProperty := NewTypeProperty()
	valueProperty.Name = "value"
	valueProperty.Type = "google.protobuf.Any"
	anyType.addProperty(valueProperty)
	domain.TypeModels[anyType.Name] = anyType
}

func (domain *Domain) sortedTypeNames() []string {
	typeNames := make([]string, 0)
	for typeName, _ := range domain.TypeModels {
		typeNames = append(typeNames, typeName)
	}
	sort.Strings(typeNames)
	return typeNames
}

func (domain *Domain) Description() string {
	typeNames := domain.sortedTypeNames()
	result := ""
	for _, typeName := range typeNames {
		result += domain.TypeModels[typeName].description()
	}
	return result
}
