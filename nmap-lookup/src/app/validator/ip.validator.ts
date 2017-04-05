import { FormControl } from '@angular/forms';


export class IpValidator {
    static ipValidator(control: FormControl): {[s:string]: boolean} {
        if (control.value && (IpValidator.verifyIPv4(<string>control.value) || IpValidator.verifyHostName(<string>control.value))) {
            return null;
        }
        return {"invalidIPv4": true}
    }

    private static verifyIPv4(input :string) : boolean {
        let result : boolean = true;
        if (input.length < 7) return false;
        if (input.charAt(0) == '.' || input.charAt(input.length-1) == '.') return false;
        let tokens: string[] = input.split('.');
        tokens.forEach((el : string)=> {
            if (el.charAt(0) == '0' && el.length > 1) {
                result = false;
            }
            let output = +el;
            if (isNaN(output)) {
                result = false;
            }
            if(output<0 || output>255) {
                result = false;
            }
        });
        return result;
    }

    private static verifyHostName(input :string) : boolean {
        if (!input.match(/^(([a-zA-Z]|[a-zA-Z][a-zA-Z\-]*[a-zA-Z])\.)*([A-Za-z]|[A-Za-z][A-Za-z\-]*[A-Za-z])$/)) {
            return false;
        }
        let tokens: string[] = input.split('.');
        if (tokens.length == 1) {
            return false;
        }
        return true;
    }
}