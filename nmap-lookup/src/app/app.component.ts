import { ResponseData } from './model/response.model';
import { LookupService } from './service/lookup.service';
import { Component } from '@angular/core';

@Component({
  selector: 'my-app',
  templateUrl: './app.component.html',
  providers: [LookupService]
})
export class AppComponent  { 
  res :ResponseData =new ResponseData();;

  public handleResult(res: ResponseData) {
    this.res = res;
  }
}
