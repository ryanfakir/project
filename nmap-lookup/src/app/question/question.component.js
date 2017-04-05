"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var ip_validator_1 = require("./../validator/ip.validator");
var response_model_1 = require("./../model/response.model");
var request_model_1 = require("./../model/request.model");
var lookup_service_1 = require("./../service/lookup.service");
var core_1 = require("@angular/core");
var forms_1 = require("@angular/forms");
var QuestionComponent = (function () {
    function QuestionComponent(lookupservice, fb) {
        this.lookupservice = lookupservice;
        this.fb = fb;
        this.host = new request_model_1.Host([]);
        this.result = new core_1.EventEmitter();
    }
    QuestionComponent.prototype.ngOnInit = function () {
        this.createForm();
        this.questionControl = new forms_1.FormControl('', forms_1.Validators.compose([forms_1.Validators.required, ip_validator_1.IpValidator.ipValidator]));
        this.questionForm.addControl('hostname', this.questionControl);
    };
    QuestionComponent.prototype.createForm = function () {
        this.questionForm = new forms_1.FormGroup({});
    };
    QuestionComponent.prototype.onSubmit = function () {
        var _this = this;
        this.extractModel();
        this.loading = true;
        this.result.emit(new response_model_1.ResponseData());
        this.lookupservice.lookup(this.host).subscribe(function (res) {
            _this.result.emit(res);
            _this.loading = false;
        });
        this.questionForm.reset();
        this.host = new request_model_1.Host([]);
    };
    QuestionComponent.prototype.extractModel = function () {
        var _this = this;
        var formModel = this.questionForm.value;
        var text = formModel.hostname;
        text.split(' ').forEach(function (el) {
            _this.host.querykeys.push(el);
        });
    };
    return QuestionComponent;
}());
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], QuestionComponent.prototype, "result", void 0);
QuestionComponent = __decorate([
    core_1.Component({
        moduleId: module.id,
        selector: 'question',
        templateUrl: 'question.component.html'
    }),
    __metadata("design:paramtypes", [lookup_service_1.LookupService, forms_1.FormBuilder])
], QuestionComponent);
exports.QuestionComponent = QuestionComponent;
//# sourceMappingURL=question.component.js.map