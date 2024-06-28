import { EvalEngine } from "./rbac";
import { UserWithRoles } from "./types";

export class DefaultEngine implements EvalEngine {
  otherCode: string;
  ruleCode: string;

  constructor() {
    this.otherCode = defaultOtherCode;
    this.ruleCode = defaultRuleCode;
  }

  SetOtherCode(code: string) {
    this.otherCode = code;
  }

  SetRuleCode(code: string) {
    this.ruleCode = code;
  }

  RunRule(user: UserWithRoles, resource: object, rule: string): boolean {
    if (rule.trim() === "") return true;

    // generate rule code
    const script = this.ruleCode.replace("%s", rule);

    const ruleFunc = new Function("user", "resource", this.otherCode + script);
    const result = ruleFunc(user, resource);
    return result;
  }
}

const defaultRuleCode = `
return %s;
`;

const defaultOtherCode = `
function listHasValue(lst, val) {
	var values = Object.values(lst);
	for(var i = 0; i < values.length; i++){
		if(values[i] === val) {
			return true;
		}
	}
	return false;
}
`;
