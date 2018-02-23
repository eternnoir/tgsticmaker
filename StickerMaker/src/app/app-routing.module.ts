import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {CreatestickersetComponent} from './createstickerset/createstickerset.component';

const routes: Routes = [
  {path: '**', component: CreatestickersetComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
