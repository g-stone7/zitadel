import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, Type } from '@angular/core';
import { AbstractControl, FormControl, FormGroup } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { take } from 'rxjs';
import {
  AddLDAPProviderRequest as AdminAddLDAPProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
  UpdateLDAPProviderRequest as AdminUpdateLDAPProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { LDAPAttributes, Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddLDAPProviderRequest as MgmtAddLDAPProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateLDAPProviderRequest as MgmtUpdateLDAPProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';

@Component({
  selector: 'cnsl-provider-ldap',
  templateUrl: './provider-ldap.component.html',
})
export class ProviderLDAPComponent {
  public showOptional: boolean = false;
  public options: Options = new Options();
  public id: string | null = '';
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  private service!: ManagementService | AdminService;

  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

  public form!: FormGroup;

  public loading: boolean = false;

  public provider?: Provider.AsObject;

  constructor(
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private _location: Location,
    private breadcrumbService: BreadcrumbService,
  ) {
    this.form = new FormGroup({
      name: new FormControl('', []),
      clientId: new FormControl('', [requiredValidator]),
      clientSecret: new FormControl('', [requiredValidator]),
      scopesList: new FormControl(['openid', 'profile', 'email'], []),
    });

    this.route.data.pipe(take(1)).subscribe((data) => {
      this.serviceType = data.serviceType;

      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          this.service = this.injector.get(ManagementService as Type<ManagementService>);

          const bread: Breadcrumb = {
            type: BreadcrumbType.ORG,
            routerLink: ['/org'],
          };

          this.breadcrumbService.setBreadcrumb([bread]);
          break;
        case PolicyComponentServiceType.ADMIN:
          this.service = this.injector.get(AdminService as Type<AdminService>);

          const iamBread = new Breadcrumb({
            type: BreadcrumbType.ORG,
            name: 'Instance',
            routerLink: ['/instance'],
          });
          this.breadcrumbService.setBreadcrumb([iamBread]);
          break;
      }

      this.id = this.route.snapshot.paramMap.get('id');
      if (this.id) {
        this.getData(this.id);
      }
    });
  }

  private getData(id: string): void {
    const req =
      this.serviceType === PolicyComponentServiceType.ADMIN
        ? new AdminGetProviderByIDRequest()
        : new MgmtGetProviderByIDRequest();
    req.setId(id);
    this.service
      .getProviderByID(req)
      .then((resp) => {
        this.provider = resp.idp;
        this.loading = false;
        if (this.provider?.config?.github) {
          this.form.patchValue(this.provider.config.github);
          this.name?.setValue(this.provider.name);
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public submitForm(): void {
    this.provider ? this.updateLDAPProvider() : this.addLDAPProvider();
  }

  public addLDAPProvider(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const req = new MgmtAddLDAPProviderRequest();

      req.setName(this.name?.value);
      req.setProviderOptions(this.options);

      this.loading = true;
      (this.service as ManagementService)
        .addLDAPProvider(req)
        .then((idp) => {
          setTimeout(() => {
            this.loading = false;
            this.close();
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
          this.loading = false;
        });
    } else if (PolicyComponentServiceType.ADMIN) {
      const req = new AdminAddLDAPProviderRequest();
      req.setName(this.name?.value);
      req.setProviderOptions(this.options);

      this.loading = true;
      (this.service as AdminService)
        .addLDAPProvider(req)
        .then((idp) => {
          setTimeout(() => {
            this.loading = false;
            this.close();
          }, 2000);
        })
        .catch((error) => {
          this.loading = false;
          this.toast.showError(error);
        });
    }
  }

  public updateLDAPProvider(): void {
    if (this.provider) {
      if (this.serviceType === PolicyComponentServiceType.MGMT) {
        const req = new MgmtUpdateLDAPProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.name?.value);
        req.setProviderOptions(this.options);

        const attr = new LDAPAttributes();
        // attr.setAvatarUrlAttribute();
        // attr.setDisplayNameAttribute();
        req.setAttributes(attr);

        // req.setBaseDn();
        // req.setBindDn();
        // req.setBindPassword();
        // req.setServersList();
        // req.setStartTls();
        // req.setTimeout();
        // req.setUserBase();
        // req.setUserFiltersList();
        // req.setUserObjectClassesList();

        this.loading = true;
        (this.service as ManagementService)
          .updateLDAPProvider(req)
          .then((idp) => {
            setTimeout(() => {
              this.loading = false;
              this.close();
            }, 2000);
          })
          .catch((error) => {
            this.toast.showError(error);
            this.loading = false;
          });
      } else if (PolicyComponentServiceType.ADMIN) {
        const req = new AdminUpdateLDAPProviderRequest();
        req.setId(this.provider.id);
        req.setName(this.name?.value);
        req.setProviderOptions(this.options);

        this.loading = true;
        (this.service as AdminService)
          .updateLDAPProvider(req)
          .then((idp) => {
            setTimeout(() => {
              this.loading = false;
              this.close();
            }, 2000);
          })
          .catch((error) => {
            this.loading = false;
            this.toast.showError(error);
          });
      }
    }
  }

  public close(): void {
    this._location.back();
  }

  public get name(): AbstractControl | null {
    return this.form.get('name');
  }

  //   public get avatarUrlAttribute(): AbstractControl | null {
  //     return this.form.get('avatarUrlAttribute');
  //   }

  //   public get displayNameAttribute(): AbstractControl | null {
  //     return this.form.get('displayNameAttribute');
  //   }

  public get baseDn(): AbstractControl | null {
    return this.form.get('baseDn');
  }

  public get bindDn(): AbstractControl | null {
    return this.form.get('bindDn');
  }

  public get bindPassword(): AbstractControl | null {
    return this.form.get('bindPassword');
  }

  public get serversList(): AbstractControl | null {
    return this.form.get('serversList');
  }

  public get startTls(): AbstractControl | null {
    return this.form.get('startTls');
  }

  public get timeout(): AbstractControl | null {
    return this.form.get('timeout');
  }

  public get userBase(): AbstractControl | null {
    return this.form.get('userBase');
  }

  public get userFiltersList(): AbstractControl | null {
    return this.form.get('userFiltersList');
  }

  public get userObjectClassesList(): AbstractControl | null {
    return this.form.get('userObjectClassesList');
  }
}
