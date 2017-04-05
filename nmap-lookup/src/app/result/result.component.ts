import { ResponseData } from './../model/response.model';
import { Component, Input } from '@angular/core';
@Component({
    moduleId: module.id,
    selector: 'result',
    templateUrl: 'result.component.html'
})
export class ResultComponent {
    @Input() result: ResponseData;
}
