"use strict";
var IpValidator = (function () {
    function IpValidator() {
    }
    IpValidator.ipValidator = function (control) {
        if (control.value && (IpValidator.verifyIPv4(control.value) || IpValidator.verifyHostName(control.value))) {
            return null;
        }
        return { "invalidIPv4": true };
    };
    IpValidator.verifyIPv4 = function (input) {
        var result = true;
        if (input.length < 7)
            return false;
        if (input.charAt(0) == '.' || input.charAt(input.length - 1) == '.')
            return false;
        var tokens = input.split('.');
        tokens.forEach(function (el) {
            if (el.charAt(0) == '0' && el.length > 1) {
                result = false;
            }
            var output = +el;
            if (isNaN(output)) {
                result = false;
            }
            if (output < 0 || output > 255) {
                result = false;
            }
        });
        return result;
    };
    IpValidator.verifyHostName = function (input) {
        if (!input.match(/^(([a-zA-Z]|[a-zA-Z][a-zA-Z\-]*[a-zA-Z])\.)*([A-Za-z]|[A-Za-z][A-Za-z\-]*[A-Za-z])$/)) {
            return false;
        }
        var tokens = input.split('.');
        if (tokens.length == 1) {
            return false;
        }
        return true;
    };
    return IpValidator;
}());
exports.IpValidator = IpValidator;
//# sourceMappingURL=ip.validator.js.map