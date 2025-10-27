import { Routes } from '@angular/router';
import { CreateUrl } from './components/create-url/create-url';
import { Home } from './components/home/home';
import { RedirectHandler } from './components/redirect-handler/redirect-handler';

export const routes: Routes = [
    {
        path: '',
        component: Home,
    },
    {
        path: 'create',
        component: CreateUrl,
    },
    {
        path: ':shortCode',
        component: RedirectHandler,
    },
];
