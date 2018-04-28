package base

import (
	"fmt"
	"strings"
	"errors"
)

func (cxt *CXProgram) GetCurrentPackage () (*CXPackage, error) {
	if cxt.CurrentPackage != nil {
		return cxt.CurrentPackage, nil
	} else {
		return nil, errors.New("current module is nil")
	}
	
}

func (cxt *CXProgram) GetCurrentStruct () (*CXStruct, error) {
	if cxt.CurrentPackage != nil &&
		cxt.CurrentPackage.CurrentStruct != nil {
		return cxt.CurrentPackage.CurrentStruct, nil
	} else {
		return nil, errors.New("current module or struct is nil")
	}
	
}

func (mod *CXPackage) GetCurrentStruct () (*CXStruct, error) {
	if mod.CurrentStruct != nil {
		return mod.CurrentStruct, nil
	} else {
		return nil, errors.New("current struct is nil")
	}
	
}

func (cxt *CXProgram) GetCurrentFunction () (*CXFunction, error) {
	if cxt.CurrentPackage != nil &&
		cxt.CurrentPackage.CurrentFunction != nil {
		return cxt.CurrentPackage.CurrentFunction, nil
	} else {
		return nil, errors.New("current module or function is nil")
	}
	
}

func (mod *CXPackage) GetCurrentFunction () (*CXFunction, error) {
	if mod.CurrentFunction != nil {
		return mod.CurrentFunction, nil
	} else {
		return nil, errors.New("current function is nil")
	}
}

func (cxt *CXProgram) GetCurrentExpression () (*CXExpression, error) {
	if cxt.CurrentPackage != nil &&
		cxt.CurrentPackage.CurrentFunction != nil &&
		cxt.CurrentPackage.CurrentFunction.CurrentExpression != nil {
		return cxt.CurrentPackage.CurrentFunction.CurrentExpression, nil
	} else {
		return nil, errors.New("current module, function or expression is nil")
	}
}

func (fn *CXFunction) GetCurrentExpression () (*CXExpression, error) {
	if fn.CurrentExpression != nil {
		return fn.CurrentExpression, nil
	} else if fn.Expressions != nil {
		return fn.Expressions[0], nil
	} else {
		return nil, errors.New("current expression is nil")
	}
}

func (cxt *CXProgram) GetCurrentDefinitions () ([]*CXDefinition, error) {
	mod, err := cxt.GetCurrentPackage()

	if err == nil {
		return mod.GetCurrentDefinitions()
	} else {
		return nil, err
	}
}

func (mod *CXPackage) GetCurrentDefinitions () ([]*CXDefinition, error) {
	return mod.GetDefinitions()
}

func (mod *CXPackage) GetDefinitions () ([]*CXDefinition, error) {
	if mod.Definitions != nil {
		return mod.Definitions, nil
	} else {
		return nil, errors.New("definitions array is nil")
	}
}

func (cxt *CXProgram) GetDefinition (name string) (*CXDefinition, error) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		var found *CXDefinition
		for _, def := range mod.Definitions {
			if def.Name == name {
				found = def
				break
			}
		}
		
		if found == nil {
			return nil, errors.New(fmt.Sprintf("GetDefinition: definition '%s' not found", name))
		} else {
			return found, nil
		}
	} else {
		return nil, err
	}
}

func (strct *CXStruct) GetFields() ([]*CXField, error) {
	if strct.Fields != nil {
		return strct.Fields, nil
	} else {
		return nil, errors.New("structure has no fields")
	}
}

func (mod *CXPackage) GetFunctions() ([]*CXFunction, error) {
	// going from map to slice
	if mod.Functions != nil {
		return mod.Functions, nil
	} else {
		return nil, errors.New("module has no functions")
	}
}

func (cxt *CXProgram) GetModule (modName string) (*CXPackage, error) {
	if cxt.Packages != nil {
		var found *CXPackage
		for _, mod := range cxt.Packages {
			if modName == mod.Name {
				found = mod
				break
			}
		}
		if found != nil {
			return found, nil
		} else {
			return nil, errors.New(fmt.Sprintf("module '%s' not found", modName))
		}
		
	} else {
		return nil, errors.New(fmt.Sprintf("module '%s' not found", modName))
	}
}

func (cxt *CXProgram) GetStruct (strctName string, modName string) (*CXStruct, error) {
	var foundMod *CXPackage
	for _, mod := range cxt.Packages {
		if modName == mod.Name {
			foundMod = mod
			break
		}
	}
	var foundStrct *CXStruct
	for _, strct := range foundMod.Structs {
		if strct.Name == strctName {
			foundStrct = strct
			break
		}
	}

	if foundStrct == nil {
		//looking in imports
		typParts := strings.Split(strctName, ".")
		
		if mod, err := cxt.GetModule(modName); err == nil {
			for _, imp := range mod.Imports {
				for _, strct := range imp.Structs {
					if strct.Name == typParts[1] {
						foundStrct = strct
						break
					}
				}
			}
		}
	}

	if foundMod != nil && foundStrct != nil {
		return foundStrct, nil
	} else {
		return nil, errors.New(fmt.Sprintf("struct '%s' not found in module '%s'", strctName, modName))
	}
}

func (mod *CXPackage) GetDefinition (defName string) (*CXDefinition, error) {
	var foundDef *CXDefinition
	for _, def := range mod.Definitions {
		if def.Name == defName {
			foundDef = def
			break
		}
	}

	if foundDef != nil {
		return foundDef, nil
	} else {
		return nil, errors.New(fmt.Sprintf("definition '%s' not found in module '%s'", defName, mod.Name))
	}
}

func (cxt *CXProgram) GetFunction (fnName string, modName string) (*CXFunction, error) {
	if _, ok := NATIVE_FUNCTIONS[fnName]; ok {
		modName = CORE_MODULE
	} else if _, ok := NATIVE_FUNCTIONS[fmt.Sprintf("%s.%s", modName, fnName)]; ok {
		fnName = fmt.Sprintf("%s.%s", modName, fnName)
		modName = CORE_MODULE
	}

	// I need to first look for the function in the current module
	// if we find modName + fnName as it is in the current module, we give that one priority
	dotFn := fmt.Sprintf("%s.%s", modName, fnName)
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		for _, fn := range mod.Functions {
			if fn.Name == dotFn {
				return fn, nil
			}
		}
	}
	
	var foundMod *CXPackage
	for _, mod := range cxt.Packages {
		if modName == mod.Name {
			foundMod = mod
			break
		}
	}
	
	var foundFn *CXFunction
	if foundMod != nil {
		for _, fn := range foundMod.Functions {
			if fn.Name == fnName {
				foundFn = fn
				break
			}
		}
	} else {
		return nil, errors.New(fmt.Sprintf("module '%s' not found", modName))
	}
	

	if foundMod != nil && foundFn != nil {
		return foundFn, nil
	} else {
		return nil, errors.New(fmt.Sprintf("function '%s' not found in module '%s'", fnName, modName))
	}
	
	// if cxt.Packages != nil && cxt.Packages[modName] != nil && cxt.Packages[modName].Functions != nil && cxt.Packages[modName].Functions[fnName] != nil {
	// 	return cxt.Packages[modName].Functions[fnName], nil
	// } else {
	// 	return nil, errors.New(fmt.Sprintf("Function '%s' not found in module '%s'", fnName, modName))
	// }
}

func (fn *CXFunction) GetExpressions () ([]*CXExpression, error) {
	if fn.Expressions != nil {
		return fn.Expressions, nil
	} else {
		return nil, errors.New("function has no expressions")
	}
}

func (fn *CXFunction) GetExpression (line int) (*CXExpression, error) {
	if fn.Expressions != nil {
		if line <= len(fn.Expressions) {
			return fn.Expressions[line], nil
		} else {
			return nil, errors.New("expression line number exceeds number of expressions in function")
		}
		
	} else {
		return nil, errors.New("function has no expressions")
	}
}

func (expr *CXExpression) GetArguments () ([]*CXArgument, error) {
	if expr.Arguments != nil {
		return expr.Arguments, nil
	} else {
		return nil, errors.New("expression has no arguments")
	}
}