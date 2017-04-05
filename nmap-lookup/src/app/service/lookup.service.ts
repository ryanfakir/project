import { Observable } from 'rxjs/Observable';
import { Host } from './../model/request.model';
import { Injectable } from "@angular/core";
import { Http, Headers, Response } from "@angular/http"
import "rxjs/Rx";
import { ResponseData } from "../model/response.model";

@Injectable()
export class LookupService {
    static URI: string = 'http://localhost:8080';

    constructor(private http: Http) {}

    public lookup(request : Host): Observable<ResponseData> {
        let headers = new Headers();
        headers.append('Content-Type', 'application/json; charset=utf-8')
        return this.http.post(LookupService.URI, request, {headers: headers}).map((res: Response) => res.json())
        .catch((error : Response) => Observable.throw(error || 'Server error'));
    }
}