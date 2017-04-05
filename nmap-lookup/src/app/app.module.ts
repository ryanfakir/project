import { ResultComponent } from './result/result.component';
import { QuestionComponent } from './question/question.component';
import { HomeComponent } from './home/home.component';
import { NgModule }      from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import {ReactiveFormsModule } from '@angular/forms';

import { AppComponent } from './app.component';
import { HttpModule } from "@angular/http";

@NgModule({
  imports:      [ BrowserModule, HttpModule, ReactiveFormsModule],
  declarations: [ AppComponent, HomeComponent, QuestionComponent, ResultComponent],
  bootstrap:    [ AppComponent ]
})
export class AppModule { }

