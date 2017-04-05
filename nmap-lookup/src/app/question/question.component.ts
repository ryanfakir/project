import { IpValidator } from './../validator/ip.validator';
import { ResponseData } from './../model/response.model';
import { Host } from './../model/request.model';
import { LookupService } from './../service/lookup.service';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { FormBuilder, FormGroup, Validators, FormControl } from "@angular/forms";
@Component({
    moduleId: module.id,
    selector: 'question',
    templateUrl: 'question.component.html'
})
export class QuestionComponent implements OnInit {
    questionForm : FormGroup;
    host: Host;
    questionControl: FormControl;
    loading: boolean;
    @Output() result: EventEmitter<ResponseData>;
    constructor(private lookupservice: LookupService, private fb: FormBuilder) {
        this.host = new Host([]);
        this.result = new EventEmitter<ResponseData>();
    }

    ngOnInit(): void {
        this.createForm();
        this.questionControl = new FormControl('', Validators.compose([Validators.required, IpValidator.ipValidator]))
        this.questionForm.addControl('hostname', this.questionControl)
    }

    public createForm() {
        this.questionForm = new FormGroup({})
    }

    public onSubmit() {
        this.extractModel()
        this.loading = true;
        this.result.emit(new ResponseData())
        this.lookupservice.lookup(this.host).subscribe((res : ResponseData) => {
            this.result.emit(res);
            this.loading = false;
        });
        this.questionForm.reset();
        this.host = new Host([]);
    }

    extractModel() {
        const formModel = this.questionForm.value;
        let text : string= formModel.hostname as string
        text.split(' ').forEach((el: string) => {
            this.host.querykeys.push(el)
        })
    }
}
